package control

import (
	"fmt"
	"log"
	"sync"
	"tcache/cache"
)

type TCache struct {
	name      string
	mainCache cache.Cache
	loader    *SingleFight
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*TCache)
)

func NewTCache(name string) *TCache {
	mu.Lock()
	defer mu.Unlock()

	// 检查是否已经存在相同的 name
	if g, exists := groups[name]; exists {
		log.Printf("TCache with name %s already exists, returning existing instance", name)
		return g
	}

	g := &TCache{
		name:      name,
		mainCache: cache.Cache{CacheBytes: 1024 * 1024 * 1024},
		loader:    &SingleFight{},
	}
	groups[name] = g
	return g
}

func GetTCache(name string) *TCache {
	mu.RLock()
	t := groups[name]
	mu.RUnlock()
	return t
}

func (g *TCache) Get(key string) (value cache.ByteView, err error) {
	byteView, err := g.loader.Do(key, func() (interface{}, error) {
		if key == "" {
			return cache.ByteView{}, fmt.Errorf("key is required")
		}

		if v, ok := g.mainCache.Get(key); ok {
			log.Println("[GeeCache] hit")
			return v, nil
		}

		return cache.ByteView{}, fmt.Errorf("key is not found")
	})
	return byteView.(cache.ByteView), err
}

func (g *TCache) Set(key string, value []byte) error {
	view, err := cache.NewByteView(value)
	if err != nil {
		return err
	}
	g.mainCache.Set(key, view)
	return nil
}
