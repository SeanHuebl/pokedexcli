package internal

import (
	"sync"
	"time"
)

// Cache is a thread-safe in-memory storage structure that holds cached entries with expiration logic.
// Entries are periodically removed based on a time interval.
type Cache struct {
	entry map[string]cacheEntry // Holds the cache data as key-value pairs.
	mu    sync.Mutex            // Mutex to synchronize access to the cache.
}

// cacheEntry represents a single cache entry with a timestamp to manage expiration.
type cacheEntry struct {
	createdAt time.Time // Timestamp when the entry was created.
	val       []byte    // Stored value for the entry.
}

// NewCache creates a new Cache instance and starts a background goroutine to periodically remove expired entries.
// The interval parameter defines how frequently expired entries are checked and removed.
func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		entry: make(map[string]cacheEntry),
	}
	ticker := time.NewTicker(interval) // Initializes a ticker to trigger expiration checks.
	go func() {
		cache.reapLoop(ticker.C) // Launches reapLoop in a goroutine to handle expired entries.
	}()
	return cache
}

// Add inserts a new entry into the cache or updates an existing one with the given key and value.
// The entry is timestamped with the current time to manage expiration.
func (c *Cache) Add(key string, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entry[key] = cacheEntry{
		createdAt: time.Now(),
		val:       value,
	}
}

// Get retrieves an entry from the cache by its key.
// It returns the stored value and a boolean indicating whether the key exists in the cache.
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	value, ok := c.entry[key]
	if !ok {
		return nil, false
	}

	return value.val, true
}

// reapLoop listens to a ticker channel and removes cache entries that have expired based on the interval set.
// This function runs in a separate goroutine to ensure expired entries are removed without blocking cache operations.
func (c *Cache) reapLoop(ch <-chan time.Time) {
	for expiration := range ch {
		c.mu.Lock()

		for key, val := range c.entry {
			// If the entry's created time is before the expiration time, remove it.
			if val.createdAt.Before(expiration) {
				delete(c.entry, key)
			}
		}
		c.mu.Unlock()
	}
}
