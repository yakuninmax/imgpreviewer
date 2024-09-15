package app

import (
	"strconv"
	"strings"

	"github.com/yakuninmax/imgpreviewer/internal/cache"
	"github.com/yakuninmax/imgpreviewer/internal/downloader"
	"github.com/yakuninmax/imgpreviewer/internal/logger"
	"github.com/yakuninmax/imgpreviewer/internal/processor"
)

type App struct {
	logger     *logger.Log
	cache      *cache.Cache
	downloader *downloader.Downloader
	processor  *processor.Processor
}

// Init app.
func New(logger *logger.Log, cache *cache.Cache, downloader *downloader.Downloader, processor *processor.Processor) *App {
	return &App{
		logger:     logger,
		cache:      cache,
		downloader: downloader,
		processor:  processor,
	}
}

// Process crop request.
func (a *App) Crop(uri string, headers map[string]string) ([]byte, error) {
	// Get request parameters.
	width, heigth, imageUrl, err := getParameters(uri)
	if err != nil {
		return nil, err
	}

	// Get image.
	img, err := a.getImage(uri, imageUrl, headers)
	if err != nil {
		return nil, err
	}

	// Crop image.
	croppedImage, err := a.processor.Crop(img, width, heigth)
	if err != nil {
		return nil, err
	}

	// Put image to cache.
	err = a.cache.Put(uri, croppedImage)
	if err != nil {
		return nil, err
	}

	return croppedImage, nil
}

// Process resize request.
func (a *App) Resize(uri string, headers map[string]string) ([]byte, error) {
	// Get request parameters.
	width, heigth, imageUrl, err := getParameters(uri)
	if err != nil {
		return nil, err
	}

	// Get image.
	img, err := a.getImage(uri, imageUrl, headers)
	if err != nil {
		return nil, err
	}

	// Resize image.
	resizedImage, err := a.processor.Resize(img, width, heigth)
	if err != nil {
		return nil, err
	}

	// Put image to cache.
	err = a.cache.Put(uri, resizedImage)
	if err != nil {
		return nil, err
	}

	return resizedImage, nil
}

// Get request parameters.
func getParameters(uri string) (int, int, string, error) {
	// Split by "/".
	subStrings := strings.SplitN(uri, "/", 3)

	// Get width.
	width, err := strconv.Atoi(subStrings[0])
	if err != nil {
		return 0, 0, "", err
	}

	// Get heigth.
	heigth, err := strconv.Atoi(subStrings[1])
	if err != nil {
		return 0, 0, "", err
	}

	return width, heigth, subStrings[2], nil
}

// Get image.
func (a *App) getImage(uri, imageUrl string, headers map[string]string) ([]byte, error) {
	// Search in cache.
	data, err := a.cache.Get(uri)
	if err != nil {
		return nil, err
	}

	// If image found in cache.
	if data != nil {
		return data, nil
	}

	// If not found in cache, download image.
	data, err = a.downloader.GetImage(imageUrl, headers)
	if err != nil {
		return nil, err
	}

	return data, nil
}
