package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	entries  map[string]cacheEntry
	mut      sync.Mutex
	interval time.Duration
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		entries:  make(map[string]cacheEntry),
		interval: interval,
	}

	go cache.reapLoop()

	return cache
}

func (c *Cache) Add(key string, val []byte) {
	c.mut.Lock()
	defer c.mut.Unlock()
	c.entries[key] = cacheEntry{createdAt: time.Now(), val: val}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mut.Lock()
	defer c.mut.Unlock()
	entry, ok := c.entries[key]
	if !ok {
		return nil, ok
	}
	return entry.val, ok
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		<-ticker.C

		c.mut.Lock()
		now := time.Now()
		for k, v := range c.entries {
			if now.Sub(v.createdAt) > c.interval {
				delete(c.entries, k)
			}
		}
		c.mut.Unlock()
	}

}
