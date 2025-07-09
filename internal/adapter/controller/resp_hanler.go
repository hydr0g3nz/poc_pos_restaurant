// Updated error handling in resp_handler.go to include menu item errors

package controller

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	errs "github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/error"
)

type successResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// SuccessResp builds a success response
func SuccessResp(c *fiber.Ctx, status int, message string, data any) error {
	return c.Status(status).JSON(successResponse{
		Status:  status,
		Message: message,
		Data:    data,
	})
}

// HandleError builds an appropriate Fiber error response based on the domain error
func HandleError(c *fiber.Ctx, err error) error {
	var statusCode int
	var message string

	switch {
	// Wallet/Transaction errors
	case errors.Is(err, errs.ErrNegativeAmount):
		statusCode = http.StatusBadRequest
		message = "Amount cannot be negative"
	case errors.Is(err, errs.ErrInsufficientBalance):
		statusCode = http.StatusBadRequest
		message = "Insufficient balance"
	case errors.Is(err, errs.ErrExpiredTransaction):
		statusCode = http.StatusBadRequest
		message = "Transaction expired"
	case errors.Is(err, errs.ErrTransactionNotVerified):
		statusCode = http.StatusBadRequest
		message = "Transaction is not in 'verified' status"
	case errors.Is(err, errs.ErrAmountExceedsLimit):
		statusCode = http.StatusBadRequest
		message = "Amount exceeds maximum limit"
	case errors.Is(err, errs.ErrInvalidPaymentMethod):
		statusCode = http.StatusBadRequest
		message = "Invalid payment method"
	case errors.Is(err, errs.ErrInvalidTransactionStatus):
		statusCode = http.StatusBadRequest
		message = "Invalid transaction status"

	// General not found errors
	case errors.Is(err, errs.ErrNotFound):
		statusCode = http.StatusNotFound
		message = "Not found"

	// Category errors
	case errors.Is(err, errs.ErrCategoryNotFound):
		statusCode = http.StatusNotFound
		message = "Category not found"
	case errors.Is(err, errs.ErrDuplicateCategoryName):
		statusCode = http.StatusConflict
		message = "Category name already exists"
	case errors.Is(err, errs.ErrInvalidCategoryType):
		statusCode = http.StatusBadRequest
		message = "Invalid category type"
	case errors.Is(err, errs.ErrCannotDeleteCategoryWithItems):
		statusCode = http.StatusConflict
		message = "Cannot delete category that has menu items"

	// Menu Item errors
	case errors.Is(err, errs.ErrMenuItemNotFound):
		statusCode = http.StatusNotFound
		message = "Menu item not found"
	case errors.Is(err, errs.ErrInvalidMenuItemPrice):
		statusCode = http.StatusBadRequest
		message = "Invalid menu item price"
	case errors.Is(err, errs.ErrMenuItemOutOfStock):
		statusCode = http.StatusConflict
		message = "Menu item is out of stock"
	case errors.Is(err, errs.ErrCannotDeleteActiveMenuItem):
		statusCode = http.StatusConflict
		message = "Cannot delete active menu item"
	case errors.Is(err, errs.ErrInvalidCategoryForMenuItem):
		statusCode = http.StatusBadRequest
		message = "Invalid category for menu item"

	// User role errors
	case errors.Is(err, errs.ErrInvalidUserRole):
		statusCode = http.StatusBadRequest
		message = "Invalid user role"

	// Order errors
	case errors.Is(err, errs.ErrInvalidOrderStatus):
		statusCode = http.StatusBadRequest
		message = "Invalid order status"
	case errors.Is(err, errs.ErrOrderNotFound):
		statusCode = http.StatusNotFound
		message = "Order not found"
	case errors.Is(err, errs.ErrOrderAlreadyClosed):
		statusCode = http.StatusConflict
		message = "Order is already closed"
	case errors.Is(err, errs.ErrOrderNotOpen):
		statusCode = http.StatusBadRequest
		message = "Order is not open"
	case errors.Is(err, errs.ErrOrderNotClosed):
		statusCode = http.StatusBadRequest
		message = "Order is not closed"
	case errors.Is(err, errs.ErrEmptyOrder):
		statusCode = http.StatusBadRequest
		message = "Order has no items"
	case errors.Is(err, errs.ErrCannotModifyClosedOrder):
		statusCode = http.StatusConflict
		message = "Cannot modify closed order"

	// Table errors
	case errors.Is(err, errs.ErrTableNotFound):
		statusCode = http.StatusNotFound
		message = "Table not found"
	case errors.Is(err, errs.ErrDuplicateTableNumber):
		statusCode = http.StatusConflict
		message = "Table number already exists"
	case errors.Is(err, errs.ErrInvalidTableNumber):
		statusCode = http.StatusBadRequest
		message = "Invalid table number"
	case errors.Is(err, errs.ErrTableAlreadyHasOpenOrder):
		statusCode = http.StatusConflict
		message = "Table already has an open order"
	case errors.Is(err, errs.ErrCannotDeleteTableWithOrders):
		statusCode = http.StatusConflict
		message = "Cannot delete table with existing orders"
	case errors.Is(err, errs.ErrTableNotAvailable):
		statusCode = http.StatusConflict
		message = "Table is not available"

	// Quantity and validation errors
	case errors.Is(err, errs.ErrInvalidQuantity):
		statusCode = http.StatusBadRequest
		message = "Invalid quantity"
	case errors.Is(err, errs.ErrInsufficientQuantity):
		statusCode = http.StatusBadRequest
		message = "Insufficient quantity"
	case errors.Is(err, errs.ErrInvalidOrderItemQuantity):
		statusCode = http.StatusBadRequest
		message = "Invalid order item quantity"

	// QR Code errors
	case errors.Is(err, errs.ErrInvalidQRCode):
		statusCode = http.StatusBadRequest
		message = "Invalid QR code"

	// Payment errors
	case errors.Is(err, errs.ErrPaymentNotFound):
		statusCode = http.StatusNotFound
		message = "Payment not found"
	case errors.Is(err, errs.ErrPaymentAlreadyExists):
		statusCode = http.StatusConflict
		message = "Payment already exists for this order"
	case errors.Is(err, errs.ErrInvalidPaymentAmount):
		statusCode = http.StatusBadRequest
		message = "Invalid payment amount"

	// Business hours and operational errors
	case errors.Is(err, errs.ErrInvalidBusinessHours):
		statusCode = http.StatusBadRequest
		message = "Operation outside business hours"
	case errors.Is(err, errs.ErrKitchenNotAvailable):
		statusCode = http.StatusServiceUnavailable
		message = "Kitchen is not available"
	case errors.Is(err, errs.ErrPrinterNotAvailable):
		statusCode = http.StatusServiceUnavailable
		message = "Printer is not available"

	// Date and time errors
	case errors.Is(err, errs.ErrInvalidDateRange):
		statusCode = http.StatusBadRequest
		message = "Invalid date range"
	case errors.Is(err, errs.ErrInvalidRevenueDate):
		statusCode = http.StatusBadRequest
		message = "Invalid revenue date"

	default:
		statusCode = http.StatusInternalServerError
		message = "Something went wrong"
	}

	return c.Status(statusCode).JSON(ErrorResponse{
		Status:  statusCode,
		Message: message,
	})
}
