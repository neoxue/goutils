package main

import (
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"time"
)

func main() {
	serverUnchange := "mctestserver1:11212"
	fmt.Println(serverUnchange)
	//create a handle
	//time.Sleep(100 * time.Second)
	a := 1

	for true {
		serverChanged := "mctestserver2:11211"
		mc := memcache.New(serverChanged)
		if mc == nil {
			fmt.Println("memcache New failed")
		}
		a++
		fmt.Println(a)
		//set key-value
		//mc.Set(&memcache.Item{Key: "foo", Value: []byte("my value" + string(a))})

		//get key's value
		it, error := mc.Get("foo")

		fmt.Println(error)
		fmt.Println(it)
		time.Sleep(1 * time.Second)
	}

}
