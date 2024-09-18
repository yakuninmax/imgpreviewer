package downloader

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var ErrInvalidFileType = errors.New("invalid file type")

type Downloader struct {
	client *http.Client
}

// Create new http client.
func New(timeout time.Duration) *Downloader {
	client := http.Client{Timeout: timeout}

	return &Downloader{&client}
}

// Get image.
func (d *Downloader) GetImage(url string, headers map[string][]string) ([]byte, error) {
	// Create request context.
	ctx := context.Background()

	// Create request.
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// Copy request headers.
	request.Header = headers

	// Send request.
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Check response status.
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error from remote server: %v", response.Status)
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
