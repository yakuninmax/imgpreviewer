package downloader

import (
	"errors"
	"io"
	"net/http"
	"time"
)

var (
	ErrInvalidFileType = errors.New("invalid file type")
)

type Downloader struct {
	client *http.Client
}

// Create new http client.
func New(timeout time.Duration) *Downloader {
	client := http.Client{Timeout: timeout}

	return &Downloader{&client}
}

// Get image.
func (d *Downloader) GetImage(url string, headers map[string]string) ([]byte, error) {
	// Create request.
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// Add headers for request.
	for key, value := range headers {
		request.Header.Add(key, value)
	}

	// Send request.
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	// Get body bytes.
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// Check if content is jpeg.
	contentType := http.DetectContentType(body)
	if contentType != "image/jpeg" {
		return nil, ErrInvalidFileType
	}

	return body, nil
}
