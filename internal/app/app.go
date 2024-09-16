package app

import (
	"fmt"
)

type logger interface {
	Info(string)
	Warn(string)
	Error(string)
	Debug(string)
}

type cache interface {
	Get(uri string) ([]byte, error)
	Put(uri string, data []byte) error
}

type downloader interface {
	GetImage(url string, headers map[string][]string) ([]byte, error)
}

type processor interface {
	Crop(data []byte, width, height int) ([]byte, error)
	Resize(data []byte, width, height int) ([]byte, error)
}

type App struct {
	logger     logger
	cache      cache
	downloader downloader
	processor  processor
}

// Init app.
func New(logg logger, cache cache, downloader downloader, processor processor) *App {
	return &App{
		logger:     logg,
		cache:      cache,
		downloader: downloader,
		processor:  processor,
	}
}

// Process crop request.
func (a *App) Crop(width, height int, url string, headers map[string][]string) ([]byte, error) {
	// Get image cache key
	cacheKey := getCacheKey(width, height, url, "crop")

	// Get image.
	img, cached, err := a.getImage(cacheKey, url, headers)
	if err != nil {
		return nil, err
	}

	// Crop image.
	croppedImage, err := a.processor.Crop(img, width, height)
	if err != nil {
		return nil, err
	}

	// Put image to cache if it not from cache.
	if !cached {
		err = a.cache.Put(cacheKey, croppedImage)
		if err != nil {
			return nil, err
		}
	}

	return croppedImage, nil
}

// Process resize request.
func (a *App) Resize(width, height int, url string, headers map[string][]string) ([]byte, error) {
	// Get image cache key
	cacheKey := getCacheKey(width, height, url, "resize")

	// Get image.
	img, cached, err := a.getImage(cacheKey, url, headers)
	if err != nil {
		return nil, err
	}

	// Resize image.
	resizedImage, err := a.processor.Resize(img, width, height)
	if err != nil {
		return nil, err
	}

	// Put image to cache if it not from cache.
	if !cached {
		err = a.cache.Put(cacheKey, resizedImage)
		if err != nil {
			return nil, err
		}
	}

	return resizedImage, nil
}

// Get image cache key.
func getCacheKey(width, height int, url, action string) string {
	return fmt.Sprintf("%d-%d-%s-%s", width, height, url, action)
}

// Get image.
func (a *App) getImage(cacheKey, url string, headers map[string][]string) ([]byte, bool, error) {
	// Search in cache.
	data, err := a.cache.Get(cacheKey)
	if err != nil {
		return nil, false, err
	}

	// If image found in cache, return it.
	if data != nil {
		return data, true, nil
	}

	// If not found in cache, download image.
	data, err = a.downloader.GetImage(url, headers)
	if err != nil {
		return nil, false, err
	}

	return data, false, nil
}
