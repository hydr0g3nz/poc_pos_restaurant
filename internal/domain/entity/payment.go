package entity

import (
	"time"

	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/vo"
)

type Payment struct {
	ID        int              `json:"id"`
	OrderID   int              `json:"order_id"`
	Amount    vo.Money         `json:"amount"`
	Method    vo.PaymentMethod `json:"method"`
	Reference string           `json:"reference,omitempty"`
	PaidAt    time.Time        `json:"paid_at"`
}

// IsValid validates payment data
func (p *Payment) IsValid() bool {
	return p.OrderID > 0 && !p.Amount.IsZero() && p.Method.Valid()
}

// NewPayment creates a new payment
func NewPayment(orderID int, amount float64, method string) (*Payment, error) {
	money, err := vo.NewMoneyFromBaht(amount)
	if err != nil {
		return nil, err
	}

	paymentMethod, err := vo.NewPaymentMethod(method)
	if err != nil {
		return nil, err
	}

	return &Payment{
		OrderID: orderID,
		Amount:  money,
		Method:  paymentMethod,
		PaidAt:  time.Now(),
	}, nil
}

// AddReference adds a reference to the payment
func (p *Payment) AddReference(reference string) {
	p.Reference = reference
}
