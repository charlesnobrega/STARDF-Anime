package anilist

import (
	"sync"
	"time"
)

// CachedEntry wraps a value with a TTL
type CachedEntry struct {
	Value     interface{}
	ExpiresAt time.Time
}

// Cache is a thread-safe in-memory cache with TTL
type Cache struct {
	mu      sync.RWMutex
	entries map[string]CachedEntry
	ttl     time.Duration
}

// NewCache creates a new cache with the given TTL
func NewCache(ttl time.Duration) *Cache {
	c := &Cache{
		entries: make(map[string]CachedEntry),
		ttl:     ttl,
	}
	go c.evictLoop()
	return c
}

// Set stores a value in the cache under the given key
func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = CachedEntry{
		Value:     value,
		ExpiresAt: time.Now().Add(c.ttl),
	}
}

// Get retrieves a value from the cache. Returns (value, true) if found and not expired.
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.entries[key]
	if !ok || time.Now().After(entry.ExpiresAt) {
		return nil, false
	}
	return entry.Value, true
}

// Delete removes an entry from the cache
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
}

// Len returns the number of entries in the cache
func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}

// evictLoop periodically cleans up expired entries
func (c *Cache) evictLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for k, v := range c.entries {
			if now.After(v.ExpiresAt) {
				delete(c.entries, k)
			}
		}
		c.mu.Unlock()
	}
}

// Global metadata cache (TTL 12h since AniList data changes rarely)
var MetaCache = NewCache(12 * time.Hour)
