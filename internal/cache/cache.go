package cache

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"sync"
)

var (
	ErrNotFound    = errors.New("file not found in cache")
	ErrFileToLarge = errors.New("file size greater than cache size")
)

type storage interface {
	Path() string
	Write(name string, data []byte) error
	Read(name string) ([]byte, error)
	Delete(name string) error
	Clean() error
}

type Cache struct {
	mu      *sync.Mutex
	size    int64
	queue   *queue
	files   map[string]*item
	storage storage
}

type file struct {
	url  string
	size int64 // image size in bytes
	name string
}

func New(size int64, storage storage) *Cache {
	mutex := &sync.Mutex{}

	return &Cache{
		mu:      mutex,
		size:    size,
		queue:   newQueue(),
		files:   make(map[string]*item),
		storage: storage,
	}
}

// Get file from cache.
func (c *Cache) Get(key string) ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if file exists.
	_, exists := c.files[key]

	if !exists {
		return nil, nil
	}

	// Read file.
	img, err := c.storage.Read(c.files[key].file.name)
	if err != nil {
		return nil, err
	}

	// Move to front.
	c.queue.moveToFront(c.files[key])
	c.files[key] = c.queue.getFront()

	return img, nil
}

// Put file to cache.
func (c *Cache) Put(key string, data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Get file size.
	size := int64(len(data))

	// Check if file size greater than cache size.
	if size > c.size {
		return ErrFileToLarge
	}

	// Get file name as hash of key (url).
	name := fmt.Sprintf("%x", sha256.Sum256([]byte(key)))

	// New cache file.
	file := file{key, size, name}

	// Check if cache space available, and cleanup.
	if c.queue.size+size > c.size {
		for {
			err := c.storage.Delete(c.queue.getBack().file.name)
			if err != nil {
				return err
			}

			delete(c.files, c.queue.getBack().file.url)
			c.queue.remove(c.queue.back)

			if c.queue.size+size <= c.size {
				break
			}
		}
	}

	// Write file.
	err := c.storage.Write(file.name, data)
	if err != nil {
		return err
	}

	// Add to queue front.
	c.queue.pushFront(file)
	c.files[key] = c.queue.getFront()

	return nil
}
