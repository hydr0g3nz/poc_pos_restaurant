package entity

import (
	"time"

	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/vo"
)

// Order represents an order domain entity
type Order struct {
	ID        int            `json:"id"`
	TableID   int            `json:"table_id"`
	Status    vo.OrderStatus `json:"status"`
	Items     []*OrderItem   `json:"items,omitempty"`
	Notes     string         `json:"notes,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	ClosedAt  *time.Time     `json:"closed_at,omitempty"`
}

// IsValid validates order data
func (o *Order) IsValid() bool {
	return o.TableID > 0 && o.Status.IsValid()
}

// NewOrder creates a new order
func NewOrder(tableID int) (*Order, error) {
	return &Order{
		TableID:   tableID,
		Status:    vo.OrderStatusOpen,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// IsOpen checks if order is open
func (o *Order) IsOpen() bool {
	return o.Status == vo.OrderStatusOpen
}

// IsClosed checks if order is closed
func (o *Order) IsClosed() bool {
	return o.Status == vo.OrderStatusClosed
}

// Close closes the order
func (o *Order) Close() {
	o.Status = vo.OrderStatusClosed
	now := time.Now()
	o.ClosedAt = &now
	o.UpdatedAt = now
}

// AddNotes adds notes to the order
func (o *Order) AddNotes(notes string) {
	o.Notes = notes
	o.UpdatedAt = time.Now()
}

// CalculateTotal calculates total amount for the order
func (o *Order) CalculateTotal() vo.Money {
	total, _ := vo.NewMoney(0)
	for _, item := range o.Items {
		itemTotal := item.CalculateSubtotal()
		total = total.Add(itemTotal)
	}
	return total
}

// GetItemCount returns total number of items in the order
func (o *Order) GetItemCount() int {
	count := 0
	for _, item := range o.Items {
		count += item.Quantity
	}
	return count
}
