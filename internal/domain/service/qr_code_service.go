package service

import (
	"context"
	"fmt"
	"time"

	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/entity"
	errs "github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/error"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/infra"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/repository"
	"github.com/hydr0g3nz/poc_pos_restuarant/utils"
)

// QRCodeService provides domain logic for QR code operations
type QRCodeService interface {
	GenerateQRCodeForOrder(ctx context.Context, tableID int) string
	ValidateQRCode(ctx context.Context, qrCode string) (*entity.Order, error)
	GenerateQRCodeImage(ctx context.Context, data string) ([]byte, error)
}

type qrCodeService struct {
	baseURL   string
	orderRepo repository.OrderRepository
	generator infra.QRCodeService
}

func NewQRCodeService(baseurl string, qrCodeImageGenerator infra.QRCodeService, orderRepo repository.OrderRepository) QRCodeService {
	return &qrCodeService{
		orderRepo: orderRepo,
		baseURL:   baseurl,
		generator: qrCodeImageGenerator,
	}
}
func genrerateQrCode(orderID int) string {
	now := time.Now().Nanosecond()
	return utils.HashSha256([]byte(fmt.Sprintf("order%dtime%d", orderID, now)))
}
func (s *qrCodeService) GenerateQRCodeForOrder(ctx context.Context, orderID int) string {
	return s.baseURL + "/order/" + genrerateQrCode(orderID)
}

func (s *qrCodeService) ValidateQRCode(ctx context.Context, qrCode string) (*entity.Order, error) {

	order, err := s.orderRepo.GetOrderByQRCode(ctx, qrCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get table: %w", err)
	}
	if order == nil {
		return nil, errs.ErrTableNotFound
	}

	return order, nil
}
func (s *qrCodeService) GenerateQRCodeImage(ctx context.Context, data string) ([]byte, error) {
	return s.generator.GenerateQRCodeImage(ctx, data)
}
