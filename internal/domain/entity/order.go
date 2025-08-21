package entity

import (
	"time"

	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/vo"
)

// Order represents an order domain entity
type Order struct {
	ID                  int              `json:"id"`
	OrderNumber         int              `json:"order_number"`
	TableID             int              `json:"table_id"`
	OrderStatus         vo.OrderStatus   `json:"status"`
	PaymentStatus       vo.PaymentStatus `json:"payment_status"`
	QRCode              string           `json:"qr_code,omitempty"` // QR code for the order
	Notes               string           `json:"notes,omitempty"`
	CreatedAt           time.Time        `json:"created_at"`
	UpdatedAt           time.Time        `json:"updated_at"`
	ClosedAt            *time.Time       `json:"closed_at,omitempty"`
	SpecialInstructions string           `json:"special_requests,omitempty"` // any special requests for the order
	Subtotal            vo.Money         `json:"subtotal,omitempty"`         // calculated subtotal for the order
	Discount            vo.Money         `json:"discount,omitempty"`         // calculated discount for the order
	TaxAmount           vo.Money         `json:"tax_amount,omitempty"`       // calculated tax for the order
	ServiceCharge       vo.Money         `json:"service_charge,omitempty"`   // calculated service charge for the order
	Total               vo.Money         `json:"total,omitempty"`            // calculated total for the order
	// extension for order items
	Items []*OrderItem `json:"items,omitempty"`
}

// IsValid validates order data
func (o *Order) IsValid() bool {
	return o.TableID > 0 && o.OrderStatus.IsValid()
}

// NewOrder creates a new order
func NewOrder(tableID int) (*Order, error) {
	return &Order{
		TableID:     tableID,
		OrderStatus: vo.OrderStatusOpen,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

// IsOpen checks if order is open
func (o *Order) IsOpen() bool {
	return o.OrderStatus == vo.OrderStatusOpen
}

// IsClosed checks if order is closed
func (o *Order) IsClosed() bool {
	return o.OrderStatus == vo.OrderStatusCompleted
}

// Close closes the order
func (o *Order) Close() {
	o.OrderStatus = vo.OrderStatusCompleted
}

// AddNotes adds notes to the order
func (o *Order) AddNotes(notes string) {
	o.Notes = notes
	o.UpdatedAt = time.Now()
}

// CalculateTotal calculates total amount for the order
func (o *Order) CalculateTotal() vo.Money {
	total, _ := vo.NewMoneyFromSatang(0)
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

func (o *Order) CalculateDiscount() vo.Money {
	// Placeholder for discount logic
	// This can be extended to apply discounts based on business rules
	return o.CalculateTotal().Multiply(0.1) // Example: 10% discount
}
func (o *Order) CalculateTax() vo.Money {
	// Placeholder for tax calculation logic
	// This can be extended to apply tax based on business rules
	return o.CalculateTotal().Multiply(0.07) // Example: 7% tax
}
