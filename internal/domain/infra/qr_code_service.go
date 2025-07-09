package infra

import "context"

// QRCodeService interface for generating and validating QR codes
type QRCodeService interface {
	// GenerateQRCode generates QR code for a table
	GenerateQRCode(ctx context.Context, tableID int) (string, error)

	// ValidateQRCode validates QR code and returns table information
	ValidateQRCode(ctx context.Context, qrCode string) (*QRCodeInfo, error)

	// GenerateQRCodeImage generates QR code as image
	GenerateQRCodeImage(ctx context.Context, data string) ([]byte, error)
}

// QRCodeInfo represents QR code information
type QRCodeInfo struct {
	TableID     int    `json:"table_id"`
	TableNumber int    `json:"table_number"`
	IsValid     bool   `json:"is_valid"`
	URL         string `json:"url"`
}
