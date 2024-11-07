package internal

import (
	"sync"
	"time"
)

type Cache struct {
	entry map[string]cacheEntry
	mu    sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		entry: make(map[string]cacheEntry),
	}
	ticker := time.NewTicker(interval)
	go func() {
		cache.reapLoop(ticker.C)
	}()
	return cache
}

func (c *Cache) Add(key string, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entry[key] = cacheEntry{
		createdAt: time.Now(),
		val:       value,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	value, ok := c.entry[key]
	if !ok {
		return nil, false
	}

	return value.val, true
}

func (c *Cache) reapLoop(ch <-chan time.Time) {
	for expiration := range ch {
		c.mu.Lock()

		for key, val := range c.entry {
			if val.createdAt.Before(expiration) {
				delete(c.entry, key)
			}
		}
		c.mu.Unlock()

	}
}
