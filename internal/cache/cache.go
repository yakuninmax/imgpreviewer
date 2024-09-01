package cache

import (
	"sync"

	"github.com/yakuninmax/imgpreviewer/internal/storage"
)

type Cache struct {
	mu      sync.Mutex
	size    int // cache current size in bytes
	list    list
	items   map[string]*listItem
	storage *storage.Storage
}

type image struct {
	uri  string
	size int // image size in bytes
	file string
}

func New(storage *storage.Storage) *Cache {
	return &Cache{
		size:    0,
		list:    list{},
		items:   map[string]*listItem{},
		storage: storage,
	}
}

func (c *Cache) Put(uri string, data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return nil
}

func (c *Cache) Get(uri string) ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return nil, nil
}
