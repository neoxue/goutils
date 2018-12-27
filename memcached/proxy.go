package memcached

import (
	"net"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

const DefaultInterval = 30
const DefaultRetryTimesLimit = 10

//Proxy is a cache proxy
type Proxy struct {
	c  *memcache.Client
	Ss *XServerList
}

//NewProxy returns the proxy
func NewProxy(server ...string) *Proxy {
	ss := &XServerList{Servers: server, TimesLimit: DefaultRetryTimesLimit, Interval: DefaultInterval}
	ss.ResolveServers()
	c := memcache.NewFromSelector(ss)
	p := &Proxy{c: c, Ss: ss}
	go p.retryfailedservers()
	return p
}
func (p *Proxy) retryfailedservers() {
	for true {
		interval := p.Ss.Interval
		time.Sleep(time.Duration(interval) * time.Second)
		for addr, status := range p.Ss.statuses {
			if status.down == 1 {
				_, err := net.DialTimeout(addr.Network(), addr.String(), memcache.DefaultTimeout)
				if err == nil {
					p.Ss.markServerUp(addr)
				} else {
					p.Ss.markServerDown(addr)
				}
			}
			if status.down == 2 {
				p.Ss.ResolveServers()
			}
		}
	}
}

func (p *Proxy) handleError(key string, err error) {
	if err == nil || err == memcache.ErrCacheMiss || err == memcache.ErrCASConflict || err == memcache.ErrNotStored || err == memcache.ErrNoStats || err == memcache.ErrMalformedKey {
		return
	}
	if err == memcache.ErrNoServers {
		p.Ss.ResolveServers()
		return
	}
	seq := p.Ss.computeServer(key, len(p.Ss.Servers))
	p.Ss.markServerDown(p.Ss.addrs[seq])
}

//FlushAll proxies client.FlushAll
func (p *Proxy) FlushAll() error {
	return p.c.FlushAll()
}

//Get proxies client.Get
func (p *Proxy) Get(key string) (item *memcache.Item, err error) {
	item, err = p.c.Get(key)
	go p.handleError(key, err)
	return
}

//Touch proxies client.Touch
func (p *Proxy) Touch(key string, seconds int32) (err error) {
	return p.c.Touch(key, seconds)
}

//GetMulti proxies client.Touch
func (p *Proxy) GetMulti(keys []string) (items map[string]*memcache.Item, err error) {
	return p.c.GetMulti(keys)
}

//Set proxies client.Set
func (p *Proxy) Set(item *memcache.Item) (err error) {
	err = p.c.Set(item)
	go p.handleError(item.Key, err)
	return
}

//Add proxies client.Add
func (p *Proxy) Add(item *memcache.Item) (err error) {
	return p.c.Add(item)
}

//Replace proxies client.Replace
func (p *Proxy) Replace(item *memcache.Item) (err error) {
	return p.c.Replace(item)
}

//CompareAndSwap proxies client.CompareAndSwap
func (p *Proxy) CompareAndSwap(item *memcache.Item) error {
	return p.c.CompareAndSwap(item)
}

//Delete proxies client.Delete
func (p *Proxy) Delete(key string) (err error) {
	return p.c.Delete(key)
}

//DeleteAll proxies client.DeleteAll
func (p *Proxy) DeleteAll() error {
	return p.c.DeleteAll()
}

//Increment proxies client.Increment
func (p *Proxy) Increment(key string, delta uint64) (newValue uint64, err error) {
	return p.c.Increment(key, delta)
}

//Decrement proxies client.Dcrement
func (p *Proxy) Decrement(key string, delta uint64) (newValue uint64, err error) {
	return p.c.Decrement(key, delta)
}
