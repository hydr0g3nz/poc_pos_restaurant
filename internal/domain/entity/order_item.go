package entity

import (
	"time"

	errs "github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/error"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/vo"
)

// OrderItem represents an order item domain entity
type OrderItem struct {
	ID              int           `json:"id"`
	OrderID         int           `json:"order_id"`
	ItemID          int           `json:"item_id"`
	Quantity        int           `json:"quantity"`
	UnitPrice       vo.Money      `json:"unit_price"`
	Name            string        `json:"name"`
	Discount        vo.Money      `json:"discount,omitempty"` // optional discount for this item
	Total           vo.Money      `json:"total"`
	SpecialReq      string        `json:"special_requests,omitempty"` // any special requests for this item
	ItemStatus      vo.ItemStatus `json:"item_status"`                // status of the item in the order
	OrderNumber     string        `json:"order_number"`               // order number for reference
	KitchenTicketID int           `json:"kitchen_id,omitempty"`
	KitchenStation  string        `json:"kitchen_station,omitempty"` // optional kitchen ID for tracking
	KitchenNotes    string        `json:"kitchen_notes,omitempty"`   // notes for the kitchen
	ServedAt        *time.Time    `json:"served_at,omitempty"`       // time when the item was served
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
}

// IsValid validates order item data
func (oi *OrderItem) IsValid() bool {
	if oi.Quantity <= 0 || oi.ItemID <= 0 || oi.OrderID <= 0 {
		return false
	}
	if oi.UnitPrice.IsZero() {
		return false
	}
	if oi.ItemStatus == "" {
		return false
	}
	return true

}

// NewOrderItem creates a new order item
func NewOrderItem(orderID, itemID int, quantity int, unitPrice float64, name string) (*OrderItem, error) {
	money, err := vo.NewMoneyFromBaht(unitPrice)
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
	subtotal, _ := vo.NewMoneyFromBaht(oi.UnitPrice.AmountBaht() * float64(oi.Quantity))
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
