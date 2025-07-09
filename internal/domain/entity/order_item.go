package entity

import (
	"time"

	errs "github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/error"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/vo"
)

// OrderItem represents an order item domain entity
type OrderItem struct {
	ID        int       `json:"id"`
	OrderID   int       `json:"order_id"`
	ItemID    int       `json:"item_id"`
	Quantity  int       `json:"quantity"`
	UnitPrice vo.Money  `json:"unit_price"`
	Notes     string    `json:"notes,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// IsValid validates order item data
func (oi *OrderItem) IsValid() bool {
	return oi.OrderID > 0 && oi.ItemID > 0 && oi.Quantity > 0 && oi.UnitPrice.Amount() >= 0
}

// NewOrderItem creates a new order item
func NewOrderItem(orderID, itemID int, quantity int, unitPrice float64) (*OrderItem, error) {
	money, err := vo.NewMoney(unitPrice)
	if err != nil {
		return nil, err
	}

	return &OrderItem{
		OrderID:   orderID,
		ItemID:    itemID,
		Quantity:  quantity,
		UnitPrice: money,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// CalculateSubtotal calculates subtotal for this order item
func (oi *OrderItem) CalculateSubtotal() vo.Money {
	subtotal, _ := vo.NewMoney(oi.UnitPrice.Amount() * float64(oi.Quantity))
	return subtotal
}

// UpdateQuantity updates the quantity of the order item
func (oi *OrderItem) UpdateQuantity(newQuantity int) error {
	if newQuantity <= 0 {
		return errs.ErrInvalidQuantity
	}

	oi.Quantity = newQuantity
	oi.UpdatedAt = time.Now()
	return nil
}

// AddNotes adds notes to the order item
func (oi *OrderItem) AddNotes(notes string) {
	oi.Notes = notes
	oi.UpdatedAt = time.Now()
}
