package memcached

import (
	"hash/crc32"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

// server selector

// ServerSelector is the interface that selects a memcache server
// as a function of the item's key.
//
// All ServerSelector implementations must be safe for concurrent use
// by multiple goroutines.
type ServerSelector interface {
	// PickServer returns the server address that a given item
	// should be shared onto.
	PickServer(key string) (net.Addr, error)
	Each(func(net.Addr) error) error
}

// XServerList is a simple ServerSelector. Its zero value is usable.
// Servers is the hostnames or original servers(which are not unix ip addr) to be resolved
type XServerList struct {
	TimesLimit int
	Interval   int64
	Servers    []string
	statuses   map[net.Addr]*serverstatus
	greenaddrs []net.Addr
	mu         sync.RWMutex
	addrs      []net.Addr
}

// serverstatus is the status of an mc instance
// status:
// 0 			green
// 1			yellow, means it failed in last 30s
// 2			red, 	dead and start to reset servers
// times:   	the times of retry
type serverstatus struct {
	down  int
	times int
}

// staticAddr caches the Network() and String() values from any net.Addr.
type staticAddr struct {
	ntw, str string
}

func newStaticAddr(a net.Addr) net.Addr {
	return &staticAddr{
		ntw: a.Network(),
		str: a.String(),
	}
}

func (s *staticAddr) Network() string { return s.ntw }
func (s *staticAddr) String() string  { return s.str }

// SetServers changes a XServerList's set of servers at runtime and is
// safe for concurrent use by multiple goroutines.
//
// Each server is given equal weight. A server is given more weight
// if it's listed multiple times.
//
// SetServers returns an error if any of the server names fail to
// resolve. No attempt is made to connect to the server. If any error
// is returned, no changes are made to the XServerList.
func (ss *XServerList) ResolveServers() error {
	servers := ss.Servers
	naddr := make([]net.Addr, len(servers))
	for i, server := range servers {
		if strings.Contains(server, "/") {
			addr, err := net.ResolveUnixAddr("unix", server)
			if err != nil {
				return err
			}
			naddr[i] = newStaticAddr(addr)
		} else {
			tcpaddr, err := net.ResolveTCPAddr("tcp", server)
			if err != nil {
				return err
			}
			naddr[i] = newStaticAddr(tcpaddr)
		}
	}

	ss.mu.Lock()
	defer ss.mu.Unlock()
	ss.addrs = naddr
	ss.greenaddrs = naddr
	ss.statuses = map[net.Addr]*serverstatus{}
	for _, addr := range naddr {
		ss.statuses[addr] = &serverstatus{times: 0, down: 0}
	}
	return nil
}

// Each iterates over each server calling the given function
func (ss *XServerList) Each(f func(net.Addr) error) error {
	ss.mu.RLock()
	defer ss.mu.RUnlock()
	for _, a := range ss.addrs {
		if err := f(a); nil != err {
			return err
		}
	}
	return nil
}

// keyBufPool returns []byte buffers for use by PickServer's call to
// crc32.ChecksumIEEE to avoid allocations. (but doesn't avoid the
// copies, which at least are bounded in size and small)
var keyBufPool = sync.Pool{
	New: func() interface{} {
		b := make([]byte, 256)
		return &b
	},
}

// PickServer should do those things
func (ss *XServerList) PickServer(key string) (net.Addr, error) {
	ss.mu.RLock()
	defer ss.mu.RUnlock()
	if len(ss.addrs) == 0 {
		return nil, memcache.ErrNoServers
	}
	if len(ss.addrs) == 1 {
		return ss.addrs[0], nil
	}
	seq := ss.computeServer(key, len(ss.Servers))
	addr := ss.addrs[seq]
	status, ok := ss.statuses[addr]
	if !ok || status.down == 0 {
		return addr, nil
	}
	if status.down >= 1 {
		// return the green server
		if len(ss.greenaddrs) < 1 {
			return nil, memcache.ErrNoServers
		}
		seq = ss.computeServer(key, len(ss.greenaddrs))
		if len(ss.greenaddrs) > int(seq) {
			return ss.greenaddrs[seq], nil
		}
		return nil, memcache.ErrNoServers
	}
	return nil, memcache.ErrNoServers
}

// Markserver down or reset servers
// if added, compute the greenaddrs
// when using ss.statuses[addr] twice, it could be probably changed, but the ss is the same;
func (ss *XServerList) markServerDown(addr net.Addr) {
	if status, ok := ss.statuses[addr]; ok {
		if status.down == 2 {
			ss.ResolveServers()
		} else {
			if status.down == 0 {
				status.down = 1
				ss.regenerateGreenAddrs()
			}
			status.times++
			if status.times >= ss.TimesLimit {
				status.down = 2
				ss.ResolveServers()
			}
		}
	}
	// if not ok,
	// which means that ss is doing ResolveServers or done another ResolveServers, the old addr is not available
	// then do nothing
	//ss.statuses[addr].times = 1
	//ss.statuses[addr].down = 1
	//ss.regenerateGreenAddrs()
}
func (ss *XServerList) markServerUp(addr net.Addr) {
	status, ok := ss.statuses[addr]
	if ok && status.down > 0 {
		status.down = 0
		status.times = 0
		ss.regenerateGreenAddrs()
	}
}
func (ss *XServerList) regenerateGreenAddrs() {
	ss.mu.RLock()
	defer ss.mu.RUnlock()
	newgreenaddrs := []net.Addr{}
	for _, addr2 := range ss.addrs {
		if ss.statuses[addr2].down == 0 {
			newgreenaddrs = append(newgreenaddrs, addr2)
		}
	}
	ss.greenaddrs = newgreenaddrs
}

// PickServer should do those things
func (ss *XServerList) computeServer(key string, num int) uint32 {
	bufp := keyBufPool.Get().(*[]byte)
	n := copy(*bufp, key)
	cs := crc32.ChecksumIEEE((*bufp)[:n])
	keyBufPool.Put(bufp)
	return cs % uint32(num)
}

// Retry failed Servers
func (ss *XServerList) retryFailedServers() {
	for true {
		interval := ss.Interval
		time.Sleep(time.Duration(interval) * time.Second)
		for addr, status := range ss.statuses {
			if status.down == 1 {
				_, err := net.DialTimeout(addr.Network(), addr.String(), memcache.DefaultTimeout)
				if err == nil {
					ss.markServerUp(addr)
				} else {
					ss.markServerDown(addr)
				}
			}
			if status.down == 2 {
				ss.ResolveServers()
			}
		}
	}
}
