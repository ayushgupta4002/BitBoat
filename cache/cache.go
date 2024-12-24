package cache

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	lock sync.RWMutex
	data map[string][]byte
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string][]byte),
	}
}

func (c *Cache) Set(key []byte, value []byte, ttl time.Duration) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.data[string(key)] = value
	ticker := time.NewTicker(ttl)
	go func() {
		<-ticker.C
		delete(c.data, string(key))
	}()
	return nil
}

func (c *Cache) Get(key []byte) ([]byte, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	val, ok := c.data[string(key)]
	// log.Println("get key", string(key), "value", string(val))
	if !ok {
		return nil, fmt.Errorf("key %s not found", key)
	}
	return val, nil
}

func (c *Cache) Has(key []byte) (bool, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	_, ok := c.data[string(key)]
	return ok, nil
}

func (c *Cache) Delete(key []byte) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.data, string(key))
	return nil
}
