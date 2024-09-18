package app

import (
	"errors"
	"fmt"
	"strconv"
)

var ErrNotEnoughParameters = errors.New("not enough parameters")

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
func (a *App) Crop(width, height, url string, headers map[string][]string) ([]byte, error) {
	// Get request parameters.
	w, h, u, err := getParameters(width, height, url)
	if err != nil {
		return nil, err
	}

	// Get image cache key
	cacheKey := getCacheKey(w, h, u, "crop")

	// Get image.
	img, cached, err := a.getImage(cacheKey, url, headers)
	if err != nil {
		return nil, err
	}

	// Crop image.
	croppedImage, err := a.processor.Crop(img, w, h)
	if err != nil {
		return nil, err
	}

	// Put image to cache if it not from cache.
	if !cached {
		err = a.cache.Put(cacheKey, croppedImage)
		if err != nil {
			return nil, err
		}

		a.logger.Debug("image " + url + " saved to cache")
	}

	return croppedImage, nil
}

// Process resize request.
func (a *App) Resize(width, height, url string, headers map[string][]string) ([]byte, error) {
	// Get request parameters.
	w, h, u, err := getParameters(width, height, url)
	if err != nil {
		return nil, err
	}

	// Get image cache key
	cacheKey := getCacheKey(w, h, u, "resize")

	// Get image.
	img, cached, err := a.getImage(cacheKey, url, headers)
	if err != nil {
		return nil, err
	}

	// Resize image.
	resizedImage, err := a.processor.Resize(img, w, h)
	if err != nil {
		return nil, err
	}

	// Put image to cache if it not from cache.
	if !cached {
		err = a.cache.Put(cacheKey, resizedImage)
		if err != nil {
			return nil, err
		}

		a.logger.Debug("image " + url + " saved to cache")
	}

	return resizedImage, nil
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
		a.logger.Debug("image " + url + " found in cache")
		return data, true, nil
	}

	a.logger.Debug("image " + url + " not found in cache, trying to download")
	// If not found in cache, download image.
	data, err = a.downloader.GetImage(url, headers)
	if err != nil {
		return nil, false, err
	}

	a.logger.Debug("image " + url + " successfully downloaded")

	return data, false, nil
}

// Get image cache key.
func getCacheKey(width, height int, url, action string) string {
	return fmt.Sprintf("%d-%d-%s-%s", width, height, url, action)
}

// Check request parameters.
func getParameters(width, height, imageURL string) (int, int, string, error) {
	// Check if parameters are not empty.
	if width == "" || height == "" || imageURL == "" {
		return 0, 0, "", ErrNotEnoughParameters
	}

	// Get width.
	w, err := strconv.Atoi(width)
	if err != nil {
		return 0, 0, "", err
	}

	// Get heigth.
	h, err := strconv.Atoi(height)
	if err != nil {
		return 0, 0, "", err
	}

	// Add scheme to url.
	u := "http://" + imageURL

	return w, h, u, nil
}
