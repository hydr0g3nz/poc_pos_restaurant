package infra

import (
	"context"
	"time"
)

// ImageService interface for handling images
type ImageService interface {
	// UploadImage uploads image and returns URL
	UploadImage(ctx context.Context, filename string, data []byte) (string, error)

	// DeleteImage deletes image
	DeleteImage(ctx context.Context, filename string) error

	// GetImageURL gets signed URL for image access
	GetImageURL(ctx context.Context, filename string, expiry time.Duration) (string, error)

	// ResizeImage resizes image to specified dimensions
	ResizeImage(ctx context.Context, data []byte, width, height int) ([]byte, error)

	// GenerateQRCodeImage generates QR code as image
	GenerateQRCodeImage(ctx context.Context, data string, size int) ([]byte, error)
}
