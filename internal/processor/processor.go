package processor

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"

	"github.com/nfnt/resize"
)

var ErrInvalidSize = errors.New("target size is larger than original")

type Processor struct{}

func New() *Processor {
	return &Processor{}
}

// Crop image.
func (p *Processor) Crop(data []byte, width, height int) ([]byte, error) {
	// Convert bytes to image.
	img, err := bytesToImage(data)
	if err != nil {
		return nil, err
	}

	// Check if destination size is larger than original.
	if width > img.Bounds().Dx() || height > img.Bounds().Dy() {
		return nil, ErrInvalidSize
	}

	// Crop image from center.
	// Set crop size.
	cropSize := image.Rect(0, 0, width, height)
	// Find starting point from center.
	cropSize = cropSize.Add(image.Point{img.Bounds().Dx()/2 - width/2, img.Bounds().Dy()/2 - height/2})
	// Crop image.
	croppedImage := img.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(cropSize)

	// Convert image to bytes.
	data, err = imageToBytes(croppedImage)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Resize image.
func (p *Processor) Resize(data []byte, width, height int) ([]byte, error) {
	// Convert bytes to image.
	img, err := bytesToImage(data)
	if err != nil {
		return nil, err
	}

	// Resize image.
	resizedImage := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)

	// Convert image to bytes.
	data, err = imageToBytes(resizedImage)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Convert image to bytes.
func imageToBytes(img image.Image) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, img, nil)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Convert bytes to image.
func bytesToImage(data []byte) (image.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	return img, nil
}
