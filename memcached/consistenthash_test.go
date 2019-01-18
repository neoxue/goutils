package memcached

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/stretchr/testify/assert"
)

const testServer1 = "localhost:11211"
const testServer2 = "localhost:11212"
const testServer3 = "localhost:11213"
const testServer2timeout = "test.sina.cn:11212"

func setupTestServers(t *testing.T) bool {
	startmcs()
	return setupTestServerOne(testServer1, t) && setupTestServerOne(testServer2, t) && setupTestServerOne(testServer3, t)
}
func startmcs() {
	fmt.Println(exec.Command("/usr/bin/memcached -m 64  -u memcache -l 127.0.0.1 -p 11211 -d").Output())
	fmt.Println(exec.Command("/usr/bin/memcached -m 64  -u memcache -l 127.0.0.1 -p 11212 -d").Output())
	fmt.Println(exec.Command("/usr/bin/memcached -m 64  -u memcache -l 127.0.0.1 -p 11213 -d").Output())
}

func setupTestServerOne(s string, t *testing.T) bool {
	var err error
	var c net.Conn
	if c, err = net.Dial("tcp", testServer1); err != nil {
		t.Skipf("skipping test; no server running at %s", testServer)
		return false
	} else {
		c.Write([]byte("flush_all\r\n"))
		c.Close()
	}
	return true
}

// 模拟11212 被删除情况
func TestLocalhosts(t *testing.T) {
	testStart(t, "/usr/bin/memcached -m 64 -u memcache -l 127.0.0.1 -p 11211 -d")
	testStart(t, "/usr/bin/memcached -m 64 -u memcache -l 127.0.0.1 -p 11212 -d")
	testStart(t, "/usr/bin/memcached -m 64 -u memcache -l 127.0.0.1 -p 11213 -d")
	if !setup(t) {
		return
	}
	testFailureTimeout(t)
	testFailureConnection(t)
	//
	//testSleep(t)
}

// to watch the connections, no problems
func TestSleep(t *testing.T) {
	var err error
	p := NewProxy(testServer1, testServer2, testServer3)
	foo := &memcache.Item{Key: "foo1", Value: []byte("fooval"), Flags: 123}
	testKill(t, "ps -ef|grep memcached|grep 11212|grep -v grep|awk '{print $2}'|xargs kill -9")
	testStart(t, "/usr/bin/memcached -m 64 -u memcache -l 127.0.0.1 -p 11212 -d")
	p.Ss.ResolveServers()
	time.Sleep(1 * time.Second)
	err = p.Set(foo)
	fmt.Println(err)
	testKill(t, "ps -ef|grep memcached|grep 11212|grep -v grep|awk '{print $2}'|xargs kill -9")
	testStart(t, "/usr/bin/memcached -m 64 -u memcache -l 127.0.0.1 -p 11212 -d")
	p.Ss.ResolveServers()
	time.Sleep(1 * time.Second)
	err = p.Set(foo)
	fmt.Println(err)
	testKill(t, "ps -ef|grep memcached|grep 11212|grep -v grep|awk '{print $2}'|xargs kill -9")
	testStart(t, "/usr/bin/memcached -m 64 -u memcache -l 127.0.0.1 -p 11212 -d")
	p.Ss.ResolveServers()
	time.Sleep(1 * time.Second)
	err = p.Set(foo)
	fmt.Println(err)
	testKill(t, "ps -ef|grep memcached|grep 11212|grep -v grep|awk '{print $2}'|xargs kill -9")
	testStart(t, "/usr/bin/memcached -m 64 -u memcache -l 127.0.0.1 -p 11212 -d")
	p.Ss.ResolveServers()
	time.Sleep(1 * time.Second)
	err = p.Set(foo)
	fmt.Println(err)
	testKill(t, "ps -ef|grep memcached|grep 11212|grep -v grep|awk '{print $2}'|xargs kill -9")
	testStart(t, "/usr/bin/memcached -m 64 -u memcache -l 127.0.0.1 -p 11212 -d")
	p.Ss.ResolveServers()
	time.Sleep(1 * time.Second)
	err = p.Set(foo)
	fmt.Println(err)
	testKill(t, "ps -ef|grep memcached|grep 11212|grep -v grep|awk '{print $2}'|xargs kill -9")
	testStart(t, "/usr/bin/memcached -m 64 -u memcache -l 127.0.0.1 -p 11212 -d")
	p.Ss.ResolveServers()
	fmt.Println(err)
	fmt.Println("in sleep, please watch the connections")
	time.Sleep(10 * time.Second)
}

func testFailureTimeout(t *testing.T) {
	p1 := NewProxy(testServer1, testServer2, testServer3)
	p1.FlushAll()
	p := NewProxy(testServer1, testServer2timeout, testServer3)
	p.Ss.TimesLimit = 3
	p.Ss.Interval = 5
	p.FlushAll()
	var foo *memcache.Item
	var err error
	var item *memcache.Item
	checkErr := func(err error, format string, args ...interface{}) {
		if err != nil {
			t.Fatalf(format, args...)
		}
	}
	//mustSet := mustSetF(t, p)
	// Set
	// test failure
	assert.Equal(t, 3, len(p.Ss.greenaddrs))
	foo = &memcache.Item{Key: "foo1", Value: []byte("fooval"), Flags: 123}
	err = p.Set(foo)
	assert.Equal(t, "memcache: connect timeout to 10.0.0.1:11212", err.Error())
	time.Sleep(100 * time.Millisecond)
	// test move to another ins
	assert.Equal(t, 2, len(p.Ss.greenaddrs))
	err = p.Set(foo)
	checkErr(err, "first set(foo): %v", err)
	item, err = p.Get("foo1")
	assert.Equal(t, 2, len(p.Ss.greenaddrs))
	assert.Equal(t, "fooval", string(item.Value))

	// test revive
	p.Ss.Servers[1] = testServer2
	time.Sleep(time.Duration(p.Ss.Interval*int64(p.Ss.TimesLimit+1)) * time.Second)
	item, err = p.Get("foo1")
	assert.Equal(t, memcache.ErrCacheMiss, err)
	assert.Equal(t, 3, len(p.Ss.greenaddrs))
	err = p.Set(foo)
	item, err = p.Get("foo1")
	assert.Equal(t, "fooval", string(item.Value))
}

func testFailureConnection(t *testing.T) {
	p := NewProxy(testServer1, testServer2, testServer3)
	p.Ss.TimesLimit = 3
	p.Ss.Interval = 5
	p.FlushAll()
	var foo *memcache.Item
	var err error
	var item *memcache.Item
	checkErr := func(err error, format string, args ...interface{}) {
		if err != nil {
			t.Fatalf(format, args...)
		}
	}
	assert.Equal(t, 3, len(p.Ss.greenaddrs))
	foo = &memcache.Item{Key: "foo1", Value: []byte("fooval"), Flags: 123}
	err = p.Set(foo)
	assert.Equal(t, nil, err)
	assert.Equal(t, 3, len(p.Ss.greenaddrs))

	testKill(t, "ps -ef|grep memcached|grep 11212|grep -v grep|awk '{print $2}'|xargs kill -9")

	// test failure
	assert.Equal(t, 3, len(p.Ss.greenaddrs))
	foo = &memcache.Item{Key: "foo1", Value: []byte("fooval"), Flags: 123}
	err = p.Set(foo)
	assert.Equal(t, "EOF", err.Error())
	time.Sleep(100 * time.Millisecond)

	// test move to another ins
	assert.Equal(t, 2, len(p.Ss.greenaddrs))
	err = p.Set(foo)
	checkErr(err, "first set(foo): %v", err)
	item, err = p.Get("foo1")
	assert.Equal(t, 2, len(p.Ss.greenaddrs))
	assert.Equal(t, "fooval", string(item.Value))

	// test revive
	testStart(t, "/usr/bin/memcached -m 64 -u memcache -l 127.0.0.1 -p 11212 -d")
	time.Sleep(time.Duration(p.Ss.Interval*int64(p.Ss.TimesLimit+2)) * time.Second)
	assert.Equal(t, 3, len(p.Ss.greenaddrs))
	item, err = p.Get("foo1")
	assert.Equal(t, memcache.ErrCacheMiss, err)
	assert.Equal(t, 3, len(p.Ss.greenaddrs))
	err = p.Set(foo)
	item, err = p.Get("foo1")
	assert.Equal(t, "fooval", string(item.Value))

}

// Run the memcached binary as a child process and connect to its unix socket.
func TestUnixSocketWithProxy(t *testing.T) {
	sock := fmt.Sprintf("/tmp/test-gomemcache-%d.sock", os.Getpid())
	cmd := exec.Command("memcached", "-s", sock)
	if err := cmd.Start(); err != nil {
		t.Skipf("skipping test; couldn't find memcached")
		return
	}
	defer cmd.Wait()
	defer cmd.Process.Kill()

	// Wait a bit for the socket to appear.
	for i := 0; i < 10; i++ {
		if _, err := os.Stat(sock); err == nil {
			break
		}
		time.Sleep(time.Duration(25*i) * time.Millisecond)
	}

	testWithClient(t, NewProxy(sock))
}

func testKill(t *testing.T, cmdstr string) {
	cmd := exec.Command("/bin/bash", "-c", cmdstr)
	bts, _ := cmd.Output()
	fmt.Println(string(bts))
}
func testStart(t *testing.T, cmdstr string) {
	cmd := exec.Command("/bin/bash", "-c", cmdstr)
	bts, _ := cmd.Output()
	fmt.Println(string(bts))
}
