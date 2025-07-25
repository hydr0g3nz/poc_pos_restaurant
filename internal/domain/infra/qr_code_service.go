package infra

import "context"

// QRCodeService interface for generating and validating QR codes
type QRCodeService interface {
	// GenerateQRCodeImage generates QR code as image
	GenerateQRCodeImage(ctx context.Context, data string) ([]byte, error)
}

// QRCodeInfo represents QR code information
