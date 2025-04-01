package cache

import (
	"sync"
)

type Cache struct {
	mu         sync.Mutex
	lru        *LRUCache
	CacheBytes int64
}

func (c *Cache) Set(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = NewLRUCache(c.CacheBytes, nil)
	}
	c.lru.Set(key, value)
}

func (c *Cache) Get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}

	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}

	return
}
