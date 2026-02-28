package cache

import (
	"reflect"
	"sync"
	"time"
)

type Cache interface {
	Get(key string, dest interface{}) bool
	Set(key string, value interface{}, duration time.Duration)
}

type item struct {
	value      interface{}
	expiration int64
}

type MemoryCache struct {
	items map[string]item
	mu    sync.RWMutex
}

func NewMemoryCache() *MemoryCache {
	c := &MemoryCache{
		items: make(map[string]item),
	}

	// Start GC goroutine to clean dead cache values every 5 minutes
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			c.mu.Lock()
			now := time.Now().UnixNano()
			for k, v := range c.items {
				if v.expiration > 0 && now > v.expiration {
					delete(c.items, k)
				}
			}
			c.mu.Unlock()
		}
	}()

	return c
}

func (c *MemoryCache) Set(key string, value interface{}, duration time.Duration) {
	var expiration int64
	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = item{
		value:      value,
		expiration: expiration,
	}
}

func (c *MemoryCache) Get(key string, dest interface{}) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	it, found := c.items[key]
	if !found {
		return false
	}

	if it.expiration > 0 && time.Now().UnixNano() > it.expiration {
		return false
	}

	// Usa reflexão para copiar o valor para o destino
	rv := reflect.ValueOf(dest)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return false
	}

	rv.Elem().Set(reflect.ValueOf(it.value))
	return true
}
