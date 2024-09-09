package downloader

import (
	"io"
	"log"
	"net/http"
	"time"
)

type Client struct {
	client *http.Client
}

// Create new http client.
func New(timeout time.Duration) *Client {
	client := http.Client{Timeout: timeout}

	return &Client{&client}
}

// Get image.
func (d *Client) GetImage(url string) {
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	println(response.Status, string(body))
}
