//go:build bughunt

package cache

import (
	"errors"
	"sync"
)

var ErrNotFound = errors.New("not found")

type Cache struct {
	mu   sync.Mutex
	data map[string]string
}

func New() *Cache {
	return &Cache{data: map[string]string{}}
}

func (c *Cache) Set(key, value string) {
	c.mu.Lock()
	c.data[key] = value
	c.mu.Unlock()
}

func (c *Cache) Get(key string) (string, error) {
	c.mu.Lock()
	v, ok := c.data[key]
	if !ok {
		return "", ErrNotFound
	}
	c.mu.Unlock()
	return v, nil
}
