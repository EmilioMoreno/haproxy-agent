package main

import (
	"sync"
	"time"
)

// CacheItem represents an item in the cache with a value and expiration time
type CacheItem struct {
	value      interface{}
	expiration time.Time
}

// Cache represents the in-memory cache
type Cache struct {
	data  map[string]CacheItem
	mutex sync.RWMutex
}

// NewCache creates a new Cache instance
func NewCache() *Cache {
	return &Cache{
		data: make(map[string]CacheItem),
	}
}

// Set adds a key-value pair to the cache with a specified expiration time
func (c *Cache) Set(key string, value interface{}, expiration time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data[key] = CacheItem{
		value:      value,
		expiration: time.Now().Add(expiration),
	}

	// Schedule a goroutine to remove the key after the specified expiration time
	go func() {
		<-time.After(expiration)
		c.mutex.Lock()
		defer c.mutex.Unlock()
		delete(c.data, key)
	}()
}

// Get retrieves the value associated with the key from the cache
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, exists := c.data[key]
	if !exists || time.Now().After(item.expiration) {
		// Key not found or expired
		return nil, false
	}

	return item.value, true
}


