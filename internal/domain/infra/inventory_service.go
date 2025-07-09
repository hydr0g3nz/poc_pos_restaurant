package infra

import "context"

// InventoryService interface for inventory management
type InventoryService interface {
	// CheckStock checks if item is in stock
	CheckStock(ctx context.Context, itemID int) (*StockInfo, error)

	// UpdateStock updates stock quantity
	UpdateStock(ctx context.Context, itemID int, quantity int) error

	// ReserveStock reserves stock for an order
	ReserveStock(ctx context.Context, orderID int, items []*StockReservation) error

	// ReleaseStock releases reserved stock
	ReleaseStock(ctx context.Context, orderID int) error

	// GetLowStockItems gets items with low stock
	GetLowStockItems(ctx context.Context) ([]*StockInfo, error)
}

// StockInfo represents stock information
type StockInfo struct {
	ItemID       int    `json:"item_id"`
	ItemName     string `json:"item_name"`
	Quantity     int    `json:"quantity"`
	MinQuantity  int    `json:"min_quantity"`
	IsLowStock   bool   `json:"is_low_stock"`
	IsOutOfStock bool   `json:"is_out_of_stock"`
}

// StockReservation represents stock reservation
type StockReservation struct {
	ItemID   int `json:"item_id"`
	Quantity int `json:"quantity"`
}
