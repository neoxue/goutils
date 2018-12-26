package memcached

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

const testServer1 = "localhost:11211"
const testServer2 = "localhost:11212"
const testServer3 = "localhost:11213"

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
	testStart11212(t)
	if !setup(t) {
		return
	}
	p := NewMemcachedProxy(testServer1, testServer2, testServer3)
	testWithClient(t, p)
	testKill11212(t)
	// test failure
	testFailure(t, p)
	fmt.Println(p.ss.greenaddrs)
	time.Sleep(30 * time.Second)
	testFailure(t, p)
	fmt.Println(p.ss.greenaddrs)
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

	testWithClient(t, NewMemcachedProxy(sock))
}

func testKill11212(t *testing.T) {
	cmd := exec.Command("/bin/bash", "-c", "ps -ef|grep memcached|grep 11212|grep -v grep|awk '{print $2}'|xargs kill -9")
	bts, _ := cmd.Output()
	fmt.Println(string(bts))
}
func testStart11212(t *testing.T) {
	cmd := exec.Command("/bin/bash", "-c", "/usr/bin/memcached -m 64 -u memcache -l 127.0.0.1 -p 11212 -d")
	bts, _ := cmd.Output()
	fmt.Println(string(bts))
}

func testFailure(t *testing.T, c *Proxy) {
	var foo *memcache.Item
	var err error
	// Set
	foo = &memcache.Item{Key: "foo", Value: []byte("fooval"), Flags: 123}
	err = c.Set(foo)
	fmt.Println("hi1")
	fmt.Println(err)
	foo = &memcache.Item{Key: "foo1", Value: []byte("fooval"), Flags: 123}
	err = c.Set(foo)
	fmt.Println("hi12")
	fmt.Println(err)
	foo = &memcache.Item{Key: "foo2", Value: []byte("fooval"), Flags: 123}
	err = c.Set(foo)
	fmt.Println("hi13")
	fmt.Println(err)
	foo = &memcache.Item{Key: "foo2", Value: []byte("fooval"), Flags: 123}
	err = c.Set(foo)
	fmt.Println("hi14")
	fmt.Println(err)
}
