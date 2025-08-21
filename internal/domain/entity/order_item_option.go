package entity

import (
	errs "github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/error"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/vo"
)

// OrderItemOption represents a menu option value for an order item
type OrderItemOption struct {
	OrderItemID     int      `json:"orderItemId"`
	OptionID        int      `json:"optionId"`
	ValueID         int      `json:"valueId"`
	AdditionalPrice vo.Money `json:"additionalPrice,omitempty"` // optional additional price for this option value
}

// IsValid checks whether the option value is valid
func (oio *OrderItemOption) IsValid() bool {
	return oio.OrderItemID > 0 && oio.OptionID > 0 && oio.ValueID > 0
}

// NewOrderItemOption creates a new order item option
func NewOrderItemOption(orderItemID int, optionID int, valueID int, additionalPrice vo.Money) (*OrderItemOption, error) {
	if orderItemID <= 0 || optionID <= 0 || valueID <= 0 {
		return nil, errs.ErrInvalidOrderItemOption
	}

	return &OrderItemOption{
		OrderItemID:     orderItemID,
		OptionID:        optionID,
		ValueID:         valueID,
		AdditionalPrice: additionalPrice,
	}, nil
}
