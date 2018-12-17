package memcached

import "github.com/bradfitz/gomemcache/memcache"

//Proxy is a cache proxy
type Proxy struct {
}

func (proxy *Proxy) Get(key string) (item *memcache.Item, err error) {

}
