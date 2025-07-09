package infra

import (
	"context"
	"time"
)

// PrinterService interface for printing receipts and kitchen orders
type PrinterService interface {
	// PrintReceipt prints customer receipt
	PrintReceipt(ctx context.Context, receipt *Receipt) error

	// PrintKitchenOrder prints kitchen order
	PrintKitchenOrder(ctx context.Context, order *KitchenOrder) error

	// PrintDailyReport prints daily sales report
	PrintDailyReport(ctx context.Context, report *DailyReport) error

	// CheckPrinterStatus checks printer status
	CheckPrinterStatus(ctx context.Context, printerID string) (*PrinterStatus, error)
}

// Receipt represents a customer receipt
type Receipt struct {
	OrderID       int            `json:"order_id"`
	TableNumber   int            `json:"table_number"`
	Items         []*ReceiptItem `json:"items"`
	Subtotal      float64        `json:"subtotal"`
	Tax           float64        `json:"tax"`
	Total         float64        `json:"total"`
	PaymentMethod string         `json:"payment_method"`
	PaidAt        time.Time      `json:"paid_at"`
	PrintedAt     time.Time      `json:"printed_at"`
}

// ReceiptItem represents an item in a receipt
type ReceiptItem struct {
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
	Subtotal  float64 `json:"subtotal"`
}

// KitchenOrder represents a kitchen order
type KitchenOrder struct {
	OrderID     int                 `json:"order_id"`
	TableNumber int                 `json:"table_number"`
	Items       []*KitchenOrderItem `json:"items"`
	Notes       string              `json:"notes"`
	CreatedAt   time.Time           `json:"created_at"`
	PrintedAt   time.Time           `json:"printed_at"`
}

// KitchenOrderItem represents an item in a kitchen order
type KitchenOrderItem struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Notes    string `json:"notes"`
}

// DailyReport represents a daily sales report
type DailyReport struct {
	Date           time.Time           `json:"date"`
	TotalRevenue   float64             `json:"total_revenue"`
	OrderCount     int                 `json:"order_count"`
	TopItems       []*TopSellingItem   `json:"top_items"`
	PaymentMethods []*PaymentMethodSum `json:"payment_methods"`
	GeneratedAt    time.Time           `json:"generated_at"`
}

// TopSellingItem represents top selling item
type TopSellingItem struct {
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Revenue  float64 `json:"revenue"`
}

// PaymentMethodSum represents payment method summary
type PaymentMethodSum struct {
	Method string  `json:"method"`
	Count  int     `json:"count"`
	Amount float64 `json:"amount"`
}

// PrinterStatus represents printer status
type PrinterStatus struct {
	PrinterID   string    `json:"printer_id"`
	IsOnline    bool      `json:"is_online"`
	PaperLevel  string    `json:"paper_level"` // "high", "medium", "low", "empty"
	LastChecked time.Time `json:"last_checked"`
	Error       string    `json:"error,omitempty"`
}
