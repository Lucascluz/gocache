// Package cache provides a thread-safe in-memory cache with TTL support.
package cache

import (
	"sync"
	"time"
)

// CacheItem represents a cached value with optional expiration.
type CacheItem struct {
	Value     any
	ExpiresAt *time.Time
}

// Cache is a thread-safe in-memory cache with TTL support.
type Cache struct {
	mu              sync.RWMutex
	items           map[string]*CacheItem
	cleanupInterval time.Duration
}

func New(cfg *Config) *Cache {
	c := &Cache{
		items: make(map[string]*CacheItem),
	}

	// Configure cleanup interval (default to 5 minutes if not provided)
	if cfg != nil && cfg.CleanupInterval > 0 {
		c.cleanupInterval = cfg.CleanupInterval
	} else {
		c.cleanupInterval = 5 * time.Minute
	}

	go c.startCleanup()

	return c
}

func (c *Cache) Set(key string, value any) error {
	return c.SetWithTTL(key, value, 0)
}

func (c *Cache) SetWithTTL(key string, value any, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	item := &CacheItem{
		Value: value,
	}

	if ttl > 0 {
		expiresAt := time.Now().Add(ttl)
		item.ExpiresAt = &expiresAt
	}

	c.items[key] = item
	return nil
}

func (c *Cache) Get(key string) (any, bool) {
	// First, read with RLock
	c.mu.RLock()
	item, exists := c.items[key]
	if !exists {
		c.mu.RUnlock()
		return nil, false
	}
	// If not expired, return quickly
	if item.ExpiresAt == nil || time.Now().Before(*item.ExpiresAt) {
		val := item.Value
		c.mu.RUnlock()
		return val, true
	}
	// Expired: upgrade to write lock to delete
	c.mu.RUnlock()
	c.mu.Lock()
	defer c.mu.Unlock()
	// Recheck under write lock
	item, exists = c.items[key]
	if !exists {
		return nil, false
	}
	if item.ExpiresAt != nil && time.Now().After(*item.ExpiresAt) {
		delete(c.items, key)
		return nil, false
	}
	return item.Value, true
}

func (c *Cache) Delete(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, exists := c.items[key]
	if exists {
		delete(c.items, key)
	}

	return exists
}

func (c *Cache) Exists(key string) bool {
	c.mu.RLock()
	item, exists := c.items[key]
	if !exists {
		c.mu.RUnlock()
		return false
	}
	if item.ExpiresAt == nil || time.Now().Before(*item.ExpiresAt) {
		c.mu.RUnlock()
		return true
	}
	// Expired: clean up under write lock
	c.mu.RUnlock()
	c.mu.Lock()
	defer c.mu.Unlock()
	item, exists = c.items[key]
	if !exists {
		return false
	}
	if item.ExpiresAt != nil && time.Now().After(*item.ExpiresAt) {
		delete(c.items, key)
		return false
	}
	return true
}

func (c *Cache) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]string, 0, len(c.items))
	now := time.Now()

	for key, item := range c.items {
		if item.ExpiresAt != nil && now.After(*item.ExpiresAt) {
			continue
		}
		keys = append(keys, key)
	}

	return keys
}

func (c *Cache) Flush() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	count := len(c.items)
	c.items = make(map[string]*CacheItem)
	return count
}

func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

func (c *Cache) startCleanup() {
	ticker := time.NewTicker(c.cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		c.cleanup()
	}
}

func (c *Cache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, item := range c.items {
		if item.ExpiresAt != nil && now.After(*item.ExpiresAt) {
			delete(c.items, key)
		}
	}
}
