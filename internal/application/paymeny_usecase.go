package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/hydr0g3nz/poc_pos_restuarant/config"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/entity"
	errs "github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/error"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/infra"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/repository"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/service"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/vo"
)

// paymentUsecase implements PaymentUsecase interface
type paymentUsecase struct {
	paymentRepo  repository.PaymentRepository
	orderRepo    repository.OrderRepository
	orderService service.OrderService
	logger       infra.Logger
	config       *config.Config
}

// NewPaymentUsecase creates a new payment usecase
func NewPaymentUsecase(
	paymentRepo repository.PaymentRepository,
	orderRepo repository.OrderRepository,
	orderService service.OrderService,
	logger infra.Logger,
	config *config.Config,
) PaymentUsecase {
	return &paymentUsecase{
		paymentRepo:  paymentRepo,
		orderRepo:    orderRepo,
		orderService: orderService,
		logger:       logger,
		config:       config,
	}
}

// ProcessPayment processes a payment for an order
func (u *paymentUsecase) ProcessPayment(ctx context.Context, req *ProcessPaymentRequest) (*PaymentResponse, error) {
	u.logger.Info("Processing payment", "orderID", req.OrderID, "amount", req.Amount, "method", req.Method)

	// Check if order exists
	order, err := u.orderRepo.GetByID(ctx, req.OrderID)
	if err != nil {
		u.logger.Error("Error getting order", "error", err, "orderID", req.OrderID)
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	if order == nil {
		u.logger.Warn("Order not found", "orderID", req.OrderID)
		return nil, errs.ErrOrderNotFound
	}

	// Check if order is closed
	if !order.IsClosed() {
		u.logger.Warn("Order is not closed", "orderID", req.OrderID)
		return nil, errs.ErrOrderNotClosed
	}

	// Check if payment already exists
	existingPayment, err := u.paymentRepo.GetByOrderID(ctx, req.OrderID)
	if err != nil {
		u.logger.Error("Error checking existing payment", "error", err, "orderID", req.OrderID)
		return nil, fmt.Errorf("failed to check existing payment: %w", err)
	}
	if existingPayment != nil {
		u.logger.Warn("Payment already exists", "orderID", req.OrderID)
		return nil, errs.ErrPaymentAlreadyExists
	}

	// Calculate order total
	total, err := u.orderService.CalculateOrderTotal(ctx, order)
	if err != nil {
		u.logger.Error("Error calculating order total", "error", err, "orderID", req.OrderID)
		return nil, fmt.Errorf("failed to calculate order total: %w", err)
	}

	// Validate payment amount
	if req.Amount != total.AmountBaht() {
		u.logger.Warn("Invalid payment amount", "orderID", req.OrderID, "expected", total.AmountBaht(), "actual", req.Amount)
		return nil, errs.ErrInvalidPaymentAmount
	}

	// Create payment entity
	payment, err := entity.NewPayment(req.OrderID, req.Amount, req.Method)
	if err != nil {
		u.logger.Error("Error creating payment entity", "error", err, "orderID", req.OrderID)
		return nil, err
	}

	// Save payment to database
	createdPayment, err := u.paymentRepo.Create(ctx, payment)
	if err != nil {
		u.logger.Error("Error creating payment", "error", err, "orderID", req.OrderID)
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	u.logger.Info("Payment processed successfully", "paymentID", createdPayment.ID, "orderID", req.OrderID)

	return u.toPaymentResponse(createdPayment), nil
}

// GetPayment retrieves payment by ID
func (u *paymentUsecase) GetPayment(ctx context.Context, id int) (*PaymentResponse, error) {
	u.logger.Debug("Getting payment", "paymentID", id)

	payment, err := u.paymentRepo.GetByID(ctx, id)
	if err != nil {
		u.logger.Error("Error getting payment", "error", err, "paymentID", id)
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}
	if payment == nil {
		u.logger.Warn("Payment not found", "paymentID", id)
		return nil, errs.ErrPaymentNotFound
	}

	return u.toPaymentResponse(payment), nil
}

// GetPaymentByOrder retrieves payment by order ID
func (u *paymentUsecase) GetPaymentByOrder(ctx context.Context, orderID int) (*PaymentResponse, error) {
	u.logger.Debug("Getting payment by order", "orderID", orderID)

	payment, err := u.paymentRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		u.logger.Error("Error getting payment by order", "error", err, "orderID", orderID)
		return nil, fmt.Errorf("failed to get payment by order: %w", err)
	}
	if payment == nil {
		u.logger.Warn("Payment not found for order", "orderID", orderID)
		return nil, errs.ErrPaymentNotFound
	}

	return u.toPaymentResponse(payment), nil
}

// ListPayments retrieves all payments with pagination
func (u *paymentUsecase) ListPayments(ctx context.Context, limit, offset int) (*PaymentListResponse, error) {
	u.logger.Debug("Listing payments", "limit", limit, "offset", offset)

	payments, err := u.paymentRepo.List(ctx, limit, offset)
	if err != nil {
		u.logger.Error("Error listing payments", "error", err)
		return nil, fmt.Errorf("failed to list payments: %w", err)
	}

	return &PaymentListResponse{
		Payments: u.toPaymentResponses(payments),
		Total:    len(payments),
		Limit:    limit,
		Offset:   offset,
	}, nil
}

// ListPaymentsByDateRange retrieves payments within date range
func (u *paymentUsecase) ListPaymentsByDateRange(ctx context.Context, startDate, endDate time.Time, limit, offset int) (*PaymentListResponse, error) {
	u.logger.Debug("Listing payments by date range", "startDate", startDate, "endDate", endDate, "limit", limit, "offset", offset)

	// Validate date range
	if startDate.After(endDate) {
		u.logger.Error("Invalid date range", "startDate", startDate, "endDate", endDate)
		return nil, errs.ErrInvalidDateRange
	}
	endDate = endDate.AddDate(0, 0, 1) // Include end date in query
	payments, err := u.paymentRepo.ListByDateRange(ctx, startDate, endDate, limit, offset)
	if err != nil {
		u.logger.Error("Error listing payments by date range", "error", err, "startDate", startDate, "endDate", endDate)
		return nil, fmt.Errorf("failed to list payments by date range: %w", err)
	}

	return &PaymentListResponse{
		Payments: u.toPaymentResponses(payments),
		Total:    len(payments),
		Limit:    limit,
		Offset:   offset,
	}, nil
}

// ListPaymentsByMethod retrieves payments by payment method
func (u *paymentUsecase) ListPaymentsByMethod(ctx context.Context, method string, limit, offset int) (*PaymentListResponse, error) {
	u.logger.Debug("Listing payments by method", "method", method, "limit", limit, "offset", offset)

	// Validate payment method
	if _, err := vo.NewPaymentMethod(method); err != nil {
		u.logger.Error("Invalid payment method", "error", err, "method", method)
		return nil, err
	}

	payments, err := u.paymentRepo.ListByMethod(ctx, method, limit, offset)
	if err != nil {
		u.logger.Error("Error listing payments by method", "error", err, "method", method)
		return nil, fmt.Errorf("failed to list payments by method: %w", err)
	}

	return &PaymentListResponse{
		Payments: u.toPaymentResponses(payments),
		Total:    len(payments),
		Limit:    limit,
		Offset:   offset,
	}, nil
}

// Helper methods for conversion

// toPaymentResponse converts entity to response
func (u *paymentUsecase) toPaymentResponse(payment *entity.Payment) *PaymentResponse {
	return &PaymentResponse{
		ID:      payment.ID,
		OrderID: payment.OrderID,
		Amount:  payment.Amount.AmountBaht(),
		Method:  payment.Method.String(),
		PaidAt:  payment.PaidAt,
	}
}

// toPaymentResponses converts slice of entities to responses
func (u *paymentUsecase) toPaymentResponses(payments []*entity.Payment) []*PaymentResponse {
	responses := make([]*PaymentResponse, len(payments))
	for i, payment := range payments {
		responses[i] = u.toPaymentResponse(payment)
	}
	return responses
}
