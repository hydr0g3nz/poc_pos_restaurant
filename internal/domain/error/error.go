// internal/domain/error/error.go (additional errors)
package errs

import (
	"time"
)

var (
	// Amount and Financial Validation
	ErrNegativeAmount        = NewValidationError("amount", "cannot be negative", nil)
	ErrAmountExceedsLimit    = NewValidationError("amount", "exceeds maximum limit", nil)
	ErrInvalidPaymentAmount  = NewValidationError("payment_amount", "must match order total", nil)
	ErrInvalidDiscountAmount = NewValidationError("discount_amount", "must be between 0 and order total", nil)

	// Order Validation
	ErrInvalidOrderStatus       = NewValidationError("order_status", "must be 'open' or 'closed'", nil)
	ErrInvalidOrderTime         = NewValidationError("order_time", "must be within business hours", nil)
	ErrInvalidOrderNote         = NewValidationError("order_note", "exceeds maximum length", nil)
	ErrInvalidOrderItemQuantity = NewValidationError("order_item_quantity", "must be greater than 0", nil)

	// Table Validation
	ErrInvalidTableNumber   = NewValidationError("table_number", "must be greater than 0", nil)
	ErrInvalidTableSeating  = NewValidationError("table_seating", "must be greater than 0", nil)
	ErrInvalidCustomerCount = NewValidationError("customer_count", "must be greater than 0", nil)

	// Menu Validation
	ErrInvalidCategoryType  = NewValidationError("category_type", "must be valid category", nil)
	ErrInvalidMenuItemPrice = NewValidationError("menu_item_price", "must be non-negative", nil)
	ErrInvalidMenuCategory  = NewValidationError("menu_category", "must be valid category", nil)

	// Quantity Validation
	ErrInvalidQuantity      = NewValidationError("quantity", "must be greater than 0", nil)
	ErrInvalidStockQuantity = NewValidationError("stock_quantity", "must be non-negative", nil)

	// Payment Method Validation
	ErrInvalidPaymentMethod = NewValidationError("payment_method", "must be 'cash', 'credit_card', or 'wallet'", nil)

	// Status Validation
	ErrInvalidTransactionStatus = NewValidationError("transaction_status", "must be valid transaction status", nil)
	ErrInvalidUserRole          = NewValidationError("user_role", "must be valid user role", nil)

	// Date and Time Validation
	ErrInvalidDateRange       = NewValidationError("date_range", "start date must be before end date", nil)
	ErrInvalidRevenueDate     = NewValidationError("revenue_date", "cannot be future date", nil)
	ErrInvalidReservationTime = NewValidationError("reservation_time", "must be in the future", nil)
	ErrInvalidWaitingTime     = NewValidationError("waiting_time", "exceeds maximum allowed time", nil)

	// QR Code and Misc Validation
	ErrInvalidQRCode         = NewValidationError("qr_code", "invalid format or expired", nil)
	ErrInvalidReceiptFormat  = NewValidationError("receipt_format", "unsupported format", nil)
	ErrInvalidTaxRate        = NewValidationError("tax_rate", "must be between 0 and 100", nil)
	ErrInvalidSpecialRequest = NewValidationError("special_request", "exceeds maximum length", nil)

	// Promo Code Validation
	ErrInvalidPromoCode     = NewValidationError("promo_code", "invalid format or not found", nil)
	ErrInvalidLoyaltyPoints = NewValidationError("loyalty_points", "must be non-negative", nil)
	// menu_option
	ErrInvalidMenuOption = NewValidationError("menu_option", "must have valid name and type", nil)
	// option_value
	ErrInvalidMenuOptionValue = NewValidationError("option_value", "must have valid name", nil)
	// payment status
	ErrInvalidPaymentStatus = NewValidationError("payment_status", "must be 'pending', 'completed', or 'failed'", nil)
	//
	ErrInvalidOrderItemOption = NewValidationError("order_item_option", "must have valid order item ID, option ID, and value ID", nil)
	//
	ErrInvalidItemStatus   = NewValidationError("item_status", "must be 'pending', 'preparing', 'ready', or 'served'", nil)
	ErrInvalidCategoryName = NewValidationError("category_name", "must be non-empty", nil)
	ErrKitchenNotFound     = NewNotFoundError("kitchen", nil)
)

// ==========================================
// Not Found Errors (404 Not Found)
// ==========================================

var (
	ErrNotFound          = NewNotFoundError("resource", nil)
	ErrOrderNotFound     = NewNotFoundError("order", nil)
	ErrTableNotFound     = NewNotFoundError("table", nil)
	ErrMenuItemNotFound  = NewNotFoundError("menu item", nil)
	ErrCategoryNotFound  = NewNotFoundError("category", nil)
	ErrPaymentNotFound   = NewNotFoundError("payment", nil)
	ErrOrderItemNotFound = NewNotFoundError("order item", nil)
)

// ==========================================
// Conflict Errors (409 Conflict)
// ==========================================

var (
	ErrDuplicateTableNumber     = NewConflictError("table", "table number already exists")
	ErrDuplicateCategoryName    = NewConflictError("category", "category name already exists")
	ErrPaymentAlreadyExists     = NewConflictError("payment", "payment already exists for this order")
	ErrTableAlreadyHasOpenOrder = NewConflictError("table", "table already has an open order")
	ErrOrderItemAlreadyExists   = NewConflictError("order item", "item already exists in order")
	ErrPromoCodeAlreadyUsed     = NewConflictError("promo code", "promo code has already been used")
)

// ==========================================
// Business Rule Errors (409 Conflict / 422 Unprocessable Entity)
// ==========================================

var (
	// Order Business Rules
	ErrOrderAlreadyClosed = NewBusinessRuleError("cannot modify closed order", map[string]interface{}{
		"rule": "order_modification",
	})
	ErrOrderNotOpen = NewBusinessRuleError("order is not in open status", map[string]interface{}{
		"rule": "order_status_check",
	})
	ErrOrderNotClosed = NewBusinessRuleError("order is not in closed status", map[string]interface{}{
		"rule": "order_status_check",
	})
	ErrEmptyOrder = NewBusinessRuleError("cannot close order without items", map[string]interface{}{
		"rule": "order_closure",
	})
	ErrCannotModifyClosedOrder = NewBusinessRuleError("cannot modify closed order", map[string]interface{}{
		"rule": "order_modification",
	})
	ErrMaxOrderItemsExceeded = NewBusinessRuleError("maximum number of order items exceeded", map[string]interface{}{
		"rule": "order_item_limit",
	})

	// Table Business Rules
	ErrTableNotAvailable = NewBusinessRuleError("table is not available for booking", map[string]interface{}{
		"rule": "table_availability",
	})
	ErrTableCapacityExceeded = NewBusinessRuleError("number of customers exceeds table capacity", map[string]interface{}{
		"rule": "table_capacity",
	})

	// Inventory Business Rules
	ErrInsufficientBalance = NewBusinessRuleError("insufficient account balance", map[string]interface{}{
		"rule": "balance_check",
	})
	ErrInsufficientQuantity = NewBusinessRuleError("insufficient quantity available", map[string]interface{}{
		"rule": "stock_availability",
	})
	ErrMenuItemOutOfStock = NewBusinessRuleError("menu item is currently out of stock", map[string]interface{}{
		"rule": "stock_availability",
	})
	ErrInsufficientLoyaltyPoints = NewBusinessRuleError("insufficient loyalty points for redemption", map[string]interface{}{
		"rule": "loyalty_points_check",
	})

	// Deletion Rules
	ErrCannotDeleteCategoryWithItems = NewBusinessRuleError("cannot delete category with existing menu items", map[string]interface{}{
		"rule": "category_deletion",
	})
	ErrCannotDeleteTableWithOrders = NewBusinessRuleError("cannot delete table with existing orders", map[string]interface{}{
		"rule": "table_deletion",
	})
	ErrCannotDeleteActiveMenuItem = NewBusinessRuleError("cannot delete active menu item", map[string]interface{}{
		"rule": "menu_item_deletion",
	})

	// Category and Menu Rules
	ErrInvalidCategoryForMenuItem = NewBusinessRuleError("menu item category does not match allowed categories", map[string]interface{}{
		"rule": "menu_item_category",
	})

	// Time-based Business Rules
	ErrExpiredTransaction = NewBusinessRuleError("transaction has expired", map[string]interface{}{
		"rule": "transaction_expiry",
	})
	ErrReservationExpired = NewBusinessRuleError("reservation has expired", map[string]interface{}{
		"rule": "reservation_expiry",
	})
	ErrPromoCodeExpired = NewBusinessRuleError("promo code has expired", map[string]interface{}{
		"rule": "promo_code_expiry",
	})
	ErrInvalidBusinessHours = NewBusinessRuleError("operation requested outside business hours", map[string]interface{}{
		"rule": "business_hours",
	})

	// Transaction Rules
	ErrTransactionNotVerified = NewBusinessRuleError("transaction is not in verified status", map[string]interface{}{
		"rule": "transaction_verification",
	})

	// Discount and Promotion Rules
	ErrDiscountNotApplicable = NewBusinessRuleError("discount is not applicable to this order", map[string]interface{}{
		"rule": "discount_applicability",
	})
	ErrServiceChargeNotApplicable = NewBusinessRuleError("service charge is not applicable", map[string]interface{}{
		"rule": "service_charge_applicability",
	})

	// Waiting List Rules
	ErrWaitingListFull = NewBusinessRuleError("waiting list has reached maximum capacity", map[string]interface{}{
		"rule": "waiting_list_capacity",
	})
)

// ==========================================
// External Service Errors (503 Service Unavailable)
// ==========================================

var (
	ErrStockNotAvailable                 = NewExternalServiceError("inventory", "stock information not available")
	ErrPrinterNotAvailable               = NewExternalServiceError("printer", "printer service unavailable")
	ErrKitchenNotAvailable               = NewExternalServiceError("kitchen", "kitchen system unavailable")
	CategoryExternal       ErrorCategory = "EXTERNAL_SERVICE"
)

// ==========================================
// Helper Functions for Context-Specific Errors
// ==========================================

// Order Errors with Context
func ErrOrderNotFoundWithID(orderID int) DomainError {
	return ErrOrderNotFound.WithField("order_id", orderID)
}

func ErrOrderAlreadyClosedWithID(orderID int) DomainError {
	return ErrOrderAlreadyClosed.WithField("order_id", orderID)
}

func ErrEmptyOrderWithID(orderID int) DomainError {
	return ErrEmptyOrder.WithField("order_id", orderID)
}

// Table Errors with Context
func ErrTableNotFoundWithID(tableID int) DomainError {
	return ErrTableNotFound.WithField("table_id", tableID)
}

func ErrTableNotFoundWithNumber(tableNumber int) DomainError {
	return ErrTableNotFound.WithField("table_number", tableNumber)
}

func ErrTableAlreadyHasOpenOrderWithContext(tableID int, openOrderID int) DomainError {
	return ErrTableAlreadyHasOpenOrder.WithDetails(map[string]interface{}{
		"table_id":      tableID,
		"open_order_id": openOrderID,
	})
}

func ErrDuplicateTableNumberWithValue(tableNumber int) DomainError {
	return ErrDuplicateTableNumber.WithField("table_number", tableNumber)
}

// Menu Item Errors with Context
func ErrMenuItemNotFoundWithID(itemID int) DomainError {
	return ErrMenuItemNotFound.WithField("item_id", itemID)
}

func ErrMenuItemOutOfStockWithID(itemID int, available int) DomainError {
	return ErrMenuItemOutOfStock.WithDetails(map[string]interface{}{
		"item_id":   itemID,
		"available": available,
	})
}

func ErrInvalidMenuItemPriceWithValue(price float64) DomainError {
	return ErrInvalidMenuItemPrice.WithField("price", price)
}

// Category Errors with Context
func ErrCategoryNotFoundWithID(categoryID int) DomainError {
	return ErrCategoryNotFound.WithField("category_id", categoryID)
}

func ErrCategoryNotFoundWithName(name string) DomainError {
	return ErrCategoryNotFound.WithField("category_name", name)
}

func ErrDuplicateCategoryNameWithValue(name string) DomainError {
	return ErrDuplicateCategoryName.WithField("category_name", name)
}

// Payment Errors with Context
func ErrPaymentNotFoundWithID(paymentID int) DomainError {
	return ErrPaymentNotFound.WithField("payment_id", paymentID)
}

func ErrPaymentNotFoundWithOrderID(orderID int) DomainError {
	return ErrPaymentNotFound.WithField("order_id", orderID)
}

func ErrPaymentAlreadyExistsWithOrderID(orderID int, existingPaymentID int) DomainError {
	return ErrPaymentAlreadyExists.WithDetails(map[string]interface{}{
		"order_id":            orderID,
		"existing_payment_id": existingPaymentID,
	})
}

func ErrInvalidPaymentAmountWithContext(expected float64, received float64) DomainError {
	return ErrInvalidPaymentAmount.WithDetails(map[string]interface{}{
		"expected_amount": expected,
		"received_amount": received,
	})
}

// Quantity Errors with Context
func ErrInvalidQuantityWithValue(quantity int) DomainError {
	return ErrInvalidQuantity.WithField("quantity", quantity)
}

func ErrInsufficientQuantityWithContext(requested int, available int) DomainError {
	return ErrInsufficientQuantity.WithDetails(map[string]interface{}{
		"requested": requested,
		"available": available,
	})
}

// Amount Errors with Context
func ErrNegativeAmountWithValue(amount float64) DomainError {
	return ErrNegativeAmount.WithField("amount", amount)
}

func ErrAmountExceedsLimitWithContext(amount float64, limit float64) DomainError {
	return ErrAmountExceedsLimit.WithDetails(map[string]interface{}{
		"amount": amount,
		"limit":  limit,
	})
}

func ErrInsufficientBalanceWithContext(balance float64, required float64) DomainError {
	return ErrInsufficientBalance.WithDetails(map[string]interface{}{
		"current_balance": balance,
		"required_amount": required,
	})
}

// Date Range Errors with Context
func ErrInvalidDateRangeWithValues(startDate time.Time, endDate time.Time) DomainError {
	return ErrInvalidDateRange.WithDetails(map[string]interface{}{
		"start_date": startDate,
		"end_date":   endDate,
	})
}

// Promo Code Errors with Context
func ErrInvalidPromoCodeWithValue(code string) DomainError {
	return ErrInvalidPromoCode.WithField("promo_code", code)
}

func ErrPromoCodeExpiredWithContext(code string, expiredAt time.Time) DomainError {
	return ErrPromoCodeExpired.WithDetails(map[string]interface{}{
		"promo_code": code,
		"expired_at": expiredAt,
	})
}

func ErrPromoCodeAlreadyUsedWithContext(code string, usedAt time.Time) DomainError {
	return ErrPromoCodeAlreadyUsed.WithDetails(map[string]interface{}{
		"promo_code": code,
		"used_at":    usedAt,
	})
}

// Stock Errors with Context
func ErrStockNotAvailableWithItem(itemID int) DomainError {
	return ErrStockNotAvailable.WithField("item_id", itemID)
}

// Business Hours Errors with Context
func ErrInvalidBusinessHoursWithTime(requestedTime time.Time) DomainError {
	return ErrInvalidBusinessHours.WithField("requested_time", requestedTime)
}

// Table Capacity Errors with Context
func ErrTableCapacityExceededWithContext(tableID int, capacity int, requested int) DomainError {
	return ErrTableCapacityExceeded.WithDetails(map[string]interface{}{
		"table_id":         tableID,
		"table_capacity":   capacity,
		"requested_guests": requested,
	})
}

// External Service Errors with Context
func ErrPrinterNotAvailableWithID(printerID string) DomainError {
	return ErrPrinterNotAvailable.WithField("printer_id", printerID)
}

func ErrKitchenNotAvailableWithReason(reason string) DomainError {
	return ErrKitchenNotAvailable.WithField("reason", reason)
}
