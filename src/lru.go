package main

import (
	"fmt"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
)

var cache = expirable.NewLRU[string, int64](1024, nil, time.Millisecond*5000)

func lruSet(key string, val int64) bool {
	return cache.Add(key, val)
}

func lruGet(key string) (int64, bool) {
	r, ok := cache.Get(key)
	if ok {
		fmt.Println("Value before expiration is found")
	}
	return r, ok
}
