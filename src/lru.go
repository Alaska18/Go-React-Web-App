package main

import (
	"sync"
	"time"
)

type CacheEntry struct {
	value  int64
	expire time.Time
}

type LRUCache struct {
	capacity int
	mutex    sync.Mutex
	cache    map[string]CacheEntry
	order    []string
	duration time.Duration
}

func newLRUCache(capacity int, expire time.Duration) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		cache:    make(map[string]CacheEntry),
		order:    make([]string, 0, capacity),
		duration: expire,
	}
}

func (c *LRUCache) set(key string, val int64) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	_, present := c.cache[key]
	evicted := false
	if present {
		c.bringFront(key)
		c.cache[key] = CacheEntry{value: val, expire: time.Now().Add(c.duration)}
		return evicted
	}
	if len(c.cache) >= c.capacity {
		c.removeLeastRecentlyUsed()
		evicted = true
	}
	c.cache[key] = CacheEntry{value: val, expire: time.Now().Add(c.duration)}
	c.addFront(key, val)
	return evicted

}

func (c *LRUCache) get(key string) (int64, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	_, present := c.cache[key]
	if present {
		return c.cache[key].value, true
	}
	return 0, false
}

func (c *LRUCache) removeFromCache(key string) {
	for i, k := range c.order {
		if k == key {
			c.order = append(c.order[:i], c.order[i+1:]...)
			break
		}
	}
}

func (c *LRUCache) bringFront(key string) {
	c.removeFromCache(key)
	c.order = append([]string{key}, c.order...)
}

func (c *LRUCache) removeLeastRecentlyUsed() {
	oldest := c.order[len(c.order)-1]
	c.removeFromCache(oldest)
}

func (c *LRUCache) addFront(key string, val int64) {
	c.order = append([]string{key}, c.order...)
}

func (c *LRUCache) removeExpired() {
	now := time.Now()
	for key, entry := range c.cache {
		if entry.expire.Before(now) {
			delete(c.cache, key)
			c.removeFromCache(key)
		}
	}
}

func (c *LRUCache) periodicCleanup() {
	ticker := time.NewTicker(c.duration)
	defer ticker.Stop()

	for range ticker.C {
		c.mutex.Lock()
		c.removeExpired()
		c.mutex.Unlock()
	}
}
