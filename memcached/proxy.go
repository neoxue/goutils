package memcached

import (
	"github.com/bradfitz/gomemcache/memcache"
)

//Proxy is a cache proxy
type Proxy struct {
	c  *memcache.Client
	ss *XServerList
}

//NewMemcachedProxy returns the proxy
func NewMemcachedProxy(server ...string) *Proxy {
	ss := &XServerList{Times: 10, Interval: 30, Servers: server}
	ss.ResolveServers()
	c := memcache.NewFromSelector(ss)
	return &Proxy{c: c, ss: ss}
}

func (p *Proxy) handleError(key string, err error) {
	if err == memcache.ErrNoServers {
		p.ss.ResolveServers()
	}
	if err == memcache.ErrServerError {
		seq := p.ss.computeServer(key, len(p.ss.Servers))
		p.ss.markServerDown(p.ss.addrs[seq])
	}
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
