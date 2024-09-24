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
func New(to time.Duration) *Downloader {
	client := http.Client{Timeout: to}

	return &Downloader{&client}
}

// Get image.
func (d *Downloader) GetImage(url string, hdr map[string][]string) ([]byte, error) {
	// Create request context.
	ctx := context.Background()

	// Create request.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// Copy request headers.
	req.Header = hdr

	// Send request.
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check response status.
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("remote server return: %v", resp.Status)
	}

	// Get body bytes.
	body, err := io.ReadAll(resp.Body)
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
