package memcached

import (
	"github.com/bradfitz/gomemcache/memcache"
)

//Proxy is a cache proxy
type Proxy struct {
	c  *memcache.Client
	ss ServerSelector
}

//NewMemcachedProxy returns the proxy
func NewMemcachedProxy(server ...string) *Proxy {
	ss := new(XServerList)
	ss.SetServers(server...)
	c := memcache.NewFromSelector(ss)
	return &Proxy{c: c, ss: ss}
}

//FlushAll proxies client.FlushAll
func (p *Proxy) FlushAll() error {
	return p.c.FlushAll()
}

//Get proxies client.Get
func (p *Proxy) Get(key string) (item *memcache.Item, err error) {
	return p.c.Get(key)
}

//Touch proxies client.Touch
func (p *Proxy) Touch(key string, seconds int32) (err error) {
	return p.c.Touch(key, seconds)
}

//GetMulti proxies client.Touch
func (p *Proxy) GetMulti(keys []string) (map[string]*memcache.Item, error) {
	return p.c.GetMulti(keys)
}

//Set proxies client.Set
func (p *Proxy) Set(item *memcache.Item) error {
	return p.c.Set(item)
}

//Add proxies client.Add
func (p *Proxy) Add(item *memcache.Item) error {
	return p.c.Add(item)
}

//Replace proxies client.Replace
func (p *Proxy) Replace(item *memcache.Item) error {
	return p.c.Replace(item)
}

//CompareAndSwap proxies client.CompareAndSwap
func (p *Proxy) CompareAndSwap(item *memcache.Item) error {
	return p.c.CompareAndSwap(item)
}

//Delete proxies client.Delete
func (p *Proxy) Delete(key string) error {
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
