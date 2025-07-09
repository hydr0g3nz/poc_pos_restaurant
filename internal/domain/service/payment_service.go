package service

import (
	"context"
	"fmt"
	"time"

	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/entity"
	errs "github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/error"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/repository"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/vo"
)

// PaymentService provides domain logic for payments
type PaymentService interface {
	// ValidatePayment validates payment before processing
	ValidatePayment(ctx context.Context, orderID int, amount vo.Money, method vo.PaymentMethod) error

	// ProcessPayment processes payment with business logic
	ProcessPayment(ctx context.Context, orderID int, amount vo.Money, method vo.PaymentMethod) (*entity.Payment, error)
}

type paymentService struct {
	paymentRepo  repository.PaymentRepository
	orderRepo    repository.OrderRepository
	orderService OrderService
}

func NewPaymentService(
	paymentRepo repository.PaymentRepository,
	orderRepo repository.OrderRepository,
	orderService OrderService,
) PaymentService {
	return &paymentService{
		paymentRepo:  paymentRepo,
		orderRepo:    orderRepo,
		orderService: orderService,
	}
}

func (s *paymentService) ValidatePayment(ctx context.Context, orderID int, amount vo.Money, method vo.PaymentMethod) error {
	// Check if order exists
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}
	if order == nil {
		return errs.ErrOrderNotFound
	}

	// Check if order is closed
	if !order.IsClosed() {
		return errs.ErrOrderNotClosed
	}

	// Check if payment already exists
	existingPayment, err := s.paymentRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to check existing payment: %w", err)
	}
	if existingPayment != nil {
		return errs.ErrPaymentAlreadyExists
	}

	// Calculate and validate amount
	total, err := s.orderService.CalculateOrderTotal(ctx, order)
	if err != nil {
		return fmt.Errorf("failed to calculate order total: %w", err)
	}

	if amount.Amount() != total.Amount() {
		return errs.ErrInvalidPaymentAmount
	}

	// Validate payment method
	if !method.Valid() {
		return errs.ErrInvalidPaymentMethod
	}

	return nil
}

func (s *paymentService) ProcessPayment(ctx context.Context, orderID int, amount vo.Money, method vo.PaymentMethod) (*entity.Payment, error) {
	// Validate payment
	if err := s.ValidatePayment(ctx, orderID, amount, method); err != nil {
		return nil, err
	}

	// Create payment
	payment := &entity.Payment{
		OrderID: orderID,
		Amount:  amount,
		Method:  method,
		PaidAt:  time.Now(),
	}

	// Save payment
	createdPayment, err := s.paymentRepo.Create(ctx, payment)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	return createdPayment, nil
}
