package cache

import (
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"sync"
	"time"

	"github.com/yakuninmax/imgpreviewer/internal/cache/storage"
)

var (
	ErrNotFound    = errors.New("file not found in cache")
	ErrFileToLarge = errors.New("file size greater than cache size")
)

type Cache struct {
	mu      *sync.Mutex
	size    int
	queue   *queue
	files   map[string]*item
	storage *storage.Storage
}

type file struct {
	uri  string
	size int // image size in bytes
	name string
}

func New(path string, size int, logger *slog.Logger) (*Cache, error) {
	var mutex = sync.Mutex{}

	// Init cache storage.
	storage, err := storage.New(path)
	if err != nil {
		return nil, err
	}

	logger.Info("cache path is " + storage.Path())

	return &Cache{
		mu:      &mutex,
		size:    size,
		queue:   newQueue(),
		files:   make(map[string]*item),
		storage: storage,
	}, nil
}

// Get file from cache.
func (c *Cache) Get(uri string) ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if file exists.
	_, exists := c.files[uri]

	if exists {
		// Read file.
		image, err := c.storage.Read(c.files[uri].data.name)
		if err != nil {
			return nil, err
		}

		// Move to front.
		c.queue.moveToFront(c.files[uri])
		c.files[uri] = c.queue.getFront()

		return image, nil
	}

	return nil, ErrNotFound
}

// Put file to cache.
func (c *Cache) Put(uri string, data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Get file size.
	size := len(data)

	// Check if file size greater than cache size.
	if size > c.size {
		return ErrFileToLarge
	}

	// Get random file name.
	name := fmt.Sprint(time.Now().Format("20060102150405"), rand.Intn(1000000))

	// New cache file.
	file := file{uri, size, name}

	// Check if cache space available, and cleanup.
	if c.queue.size+size > c.size {
		for ok := true; ok; ok = c.queue.size <= c.size {
			delete(c.files, c.queue.getBack().data.uri)
			c.queue.remove(c.queue.back)
		}
	}

	// Write file.
	err := c.storage.Write(file.name, data)
	if err != nil {
		return err
	}

	// Add to queue front.
	c.queue.pushFront(file)
	c.files[uri] = c.queue.getFront()

	return nil
}
