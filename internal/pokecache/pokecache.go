package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	val       []byte
	createdAt time.Time
}

type Cache struct {
	entries map[string]cacheEntry
	mux     *sync.Mutex
}

func NewCache(interval time.Duration) Cache {

	newCache := Cache{
		entries: make(map[string]cacheEntry),
		mux:     &sync.Mutex{},
	}
	go newCache.reapLoop(interval)
	return newCache
}

func (c *Cache) Add(key string, val []byte) {

	c.mux.Lock()
	defer c.mux.Unlock()
	c.entries[key] = cacheEntry{
		val:       val,
		createdAt: time.Now().UTC(),
	}
}

func (c *Cache) Get(key string) (data []byte, exists bool) {

	c.mux.Lock()
	defer c.mux.Unlock()
	entry, exists := c.entries[key]
	return entry.val, exists
}

func (c *Cache) reapLoop(interval time.Duration) {

	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.reap(interval)
	}
}

func (c *Cache) reap(interval time.Duration) {

	c.mux.Lock()
	defer c.mux.Unlock()
	intervalTimeAgo := time.Now().UTC().Add(-interval)
	for k, entry := range c.entries {
		if entry.createdAt.Before(intervalTimeAgo) {
			delete(c.entries, k)
		}
	}
}
