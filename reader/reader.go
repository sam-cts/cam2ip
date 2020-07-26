// Package reader.
package reader

import (
	"image"
)

// ImageReader interface
type ImageReader interface {
	// Read reads next frame from camera/video and returns image.
	Read() (img image.Image, err error)

	GetImage() (img image.Image)

	// Close closes camera/video.
	Close() error
}
