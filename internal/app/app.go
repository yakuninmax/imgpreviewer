package app

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"strconv"

	"github.com/disintegration/imaging"
)

var (
	ErrNotEnoughParameters = errors.New("not enough parameters")
	ErrInvalidSize         = errors.New("target size is larger than original")
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

type App struct {
	logger     logger
	cache      cache
	downloader downloader
}

// Init app.
func New(logg logger, cache cache, downloader downloader) *App {
	return &App{
		logger:     logg,
		cache:      cache,
		downloader: downloader,
	}
}

// Process resize request.
func (a *App) Fill(width, height, url string, headers map[string][]string) ([]byte, error) {
	// Get request parameters.
	w, h, u, err := getParameters(width, height, url)
	if err != nil {
		return nil, err
	}

	// Get image cache key
	cacheKey := getCacheKey(w, h, u, "resize")

	// Get image.
	b, cached, err := a.getImage(cacheKey, u, headers)
	if err != nil {
		return nil, err
	}

	// Bytes to image.
	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	// Check if destination size is larger than original.
	if w > img.Bounds().Dx() || h > img.Bounds().Dy() {
		return nil, ErrInvalidSize
	}

	// Resize image.
	img = imaging.Fill(img, w, h, imaging.Center, imaging.Lanczos)

	// Image to bytes.
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, img, nil)
	if err != nil {
		return nil, err
	}

	// Put image to cache if it not from cache.
	if !cached {
		err = a.cache.Put(cacheKey, buf.Bytes())
		if err != nil {
			return nil, err
		}

		a.logger.Debug("image " + url + " saved to cache")
	}

	return buf.Bytes(), nil
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
