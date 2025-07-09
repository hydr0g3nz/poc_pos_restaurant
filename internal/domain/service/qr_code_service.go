package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/entity"
	errs "github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/error"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/repository"
)

// QRCodeService provides domain logic for QR code operations
type QRCodeService interface {
	// GenerateQRCodeForTable generates QR code for table
	GenerateQRCodeForTable(ctx context.Context, tableID int) (string, error)

	// ValidateQRCode validates QR code and returns table info
	ValidateQRCode(ctx context.Context, qrCode string) (*entity.Table, error)

	// ParseTableFromQRCode parses table ID from QR code
	ParseTableFromQRCode(ctx context.Context, qrCode string) (int, error)
}

type qrCodeService struct {
	tableRepo repository.TableRepository
}

func NewQRCodeService(tableRepo repository.TableRepository) QRCodeService {
	return &qrCodeService{
		tableRepo: tableRepo,
	}
}

func (s *qrCodeService) GenerateQRCodeForTable(ctx context.Context, tableID int) (string, error) {
	// Check if table exists
	table, err := s.tableRepo.GetByID(ctx, tableID)
	if err != nil {
		return "", fmt.Errorf("failed to get table: %w", err)
	}
	if table == nil {
		return "", errs.ErrTableNotFound
	}

	// Generate QR code URL
	qrCode := fmt.Sprintf("/order?table=%d", table.TableNumber)

	// Update table with QR code
	table.QRCode = qrCode
	_, err = s.tableRepo.Update(ctx, table)
	if err != nil {
		return "", fmt.Errorf("failed to update table with QR code: %w", err)
	}

	return qrCode, nil
}

func (s *qrCodeService) ValidateQRCode(ctx context.Context, qrCode string) (*entity.Table, error) {
	// Parse table number from QR code
	tableNumber, err := s.ParseTableFromQRCode(ctx, qrCode)
	if err != nil {
		return nil, err
	}

	// Get table by number
	table, err := s.tableRepo.GetByNumber(ctx, tableNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get table: %w", err)
	}
	if table == nil {
		return nil, errs.ErrTableNotFound
	}

	return table, nil
}

func (s *qrCodeService) ParseTableFromQRCode(ctx context.Context, qrCode string) (int, error) {
	// Parse QR code format: /order?table=X
	if !strings.HasPrefix(qrCode, "/order?table=") {
		return 0, errs.ErrInvalidQRCode
	}

	// Extract table number
	tableStr := strings.TrimPrefix(qrCode, "/order?table=")
	tableNumber, err := strconv.Atoi(tableStr)
	if err != nil {
		return 0, errs.ErrInvalidQRCode
	}

	if tableNumber <= 0 {
		return 0, errs.ErrInvalidQRCode
	}

	return tableNumber, nil
}
