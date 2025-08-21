package vo

import (
	"strings"

	errs "github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/error"
)

type OrderStatus string

const (
	OrderStatusOpen      OrderStatus = "open"
	OrderStatusOrdered   OrderStatus = "ordered"
	OrderStatusCompleted OrderStatus = "completed"
	OrderCancelled       OrderStatus = "cancelled"
)

func (s OrderStatus) IsValid() bool {
	switch s {
	case OrderStatusOpen, OrderStatusOrdered, OrderStatusCompleted, OrderCancelled:
		return true
	default:
		return false
	}
}

func NewOrderStatus(status string) (OrderStatus, error) {
	s := OrderStatus(strings.ToLower(status))
	if !s.IsValid() {
		return "", errs.ErrInvalidOrderStatus
	}
	return s, nil
}

func (s OrderStatus) String() string {
	return string(s)
}
