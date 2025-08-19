package controller

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/adapter/dto"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/adapter/presenter"
	usecase "github.com/hydr0g3nz/poc_pos_restuarant/internal/application"
)

// PaymentController handles HTTP requests related to payment operations
type PaymentController struct {
	paymentUseCase usecase.PaymentUsecase
	errorPresenter presenter.ErrorPresenter
}

// NewPaymentController creates a new instance of PaymentController
func NewPaymentController(paymentUseCase usecase.PaymentUsecase, errorPresenter presenter.ErrorPresenter) *PaymentController {
	return &PaymentController{
		paymentUseCase: paymentUseCase,
		errorPresenter: errorPresenter,
	}
}

// ProcessPayment handles payment processing
func (c *PaymentController) ProcessPayment(ctx *fiber.Ctx) error {
	var req dto.ProcessPaymentRequest
	if err := ctx.BodyParser(&req); err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	if req.OrderID <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Order ID is required and must be greater than 0",
		})
	}

	if req.Amount <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Amount is required and must be greater than 0",
		})
	}

	if req.Method == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Payment method is required",
		})
	}

	response, err := c.paymentUseCase.ProcessPayment(ctx.Context(), &usecase.ProcessPaymentRequest{
		OrderID: req.OrderID,
		Amount:  req.Amount,
		Method:  req.Method,
	})
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusCreated, "Payment processed successfully", response)
}

// GetPayment handles getting payment by ID
func (c *PaymentController) GetPayment(ctx *fiber.Ctx) error {
	paymentIDParam := ctx.Params("id")
	if paymentIDParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Payment ID is required",
		})
	}

	paymentID, err := strconv.Atoi(paymentIDParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid Payment ID format",
		})
	}

	response, err := c.paymentUseCase.GetPayment(ctx.Context(), paymentID)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Payment retrieved successfully", response)
}

// GetPaymentByOrder handles getting payment by order ID
func (c *PaymentController) GetPaymentByOrder(ctx *fiber.Ctx) error {
	orderIDParam := ctx.Params("orderId")
	if orderIDParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Order ID is required",
		})
	}

	orderID, err := strconv.Atoi(orderIDParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid Order ID format",
		})
	}

	response, err := c.paymentUseCase.GetPaymentByOrder(ctx.Context(), orderID)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Payment retrieved successfully", response)
}

// ListPayments handles getting all payments
func (c *PaymentController) ListPayments(ctx *fiber.Ctx) error {
	// Parse pagination parameters
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	offset, _ := strconv.Atoi(ctx.Query("offset", "0"))

	// Validate pagination parameters
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	response, err := c.paymentUseCase.ListPayments(ctx.Context(), limit, offset)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Payments retrieved successfully", response)
}

// ListPaymentsByDateRange handles getting payments by date range
func (c *PaymentController) ListPaymentsByDateRange(ctx *fiber.Ctx) error {
	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "start_date and end_date query parameters are required",
		})
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid start_date format. Use YYYY-MM-DD",
		})
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid end_date format. Use YYYY-MM-DD",
		})
	}

	// Parse pagination parameters
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	offset, _ := strconv.Atoi(ctx.Query("offset", "0"))

	// Validate pagination parameters
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	response, err := c.paymentUseCase.ListPaymentsByDateRange(ctx.Context(), startDate, endDate, limit, offset)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Payments by date range retrieved successfully", response)
}

// ListPaymentsByMethod handles getting payments by payment method
func (c *PaymentController) ListPaymentsByMethod(ctx *fiber.Ctx) error {
	method := ctx.Query("method")
	if method == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Method query parameter is required",
		})
	}

	// Parse pagination parameters
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	offset, _ := strconv.Atoi(ctx.Query("offset", "0"))

	// Validate pagination parameters
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	response, err := c.paymentUseCase.ListPaymentsByMethod(ctx.Context(), method, limit, offset)
	if err != nil {
		return HandleError(ctx, err, c.errorPresenter)
	}

	return SuccessResp(ctx, fiber.StatusOK, "Payments by method retrieved successfully", response)
}
