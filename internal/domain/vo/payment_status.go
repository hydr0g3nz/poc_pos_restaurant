package vo

import (
	"strings"

	errs "github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/error"
)

type PaymentStatus string

const (
	PaymentStatusUnpaid PaymentStatus = "unpaid"
	// PaymentStatusPartial  PaymentStatus = "partial"
	PaymentStatusPaid     PaymentStatus = "paid"
	PaymentStatusRefunded PaymentStatus = "refunded"
)

func (s PaymentStatus) IsValid() bool {
	switch s {
	case PaymentStatusUnpaid, PaymentStatusPaid, PaymentStatusRefunded:
		return true
	default:
		return false
	}
}

func NewPaymentStatus(status string) (PaymentStatus, error) {
	s := PaymentStatus(strings.ToLower(status))
	if !s.IsValid() {
		return "", errs.ErrInvalidPaymentStatus
	}
	return s, nil
}

func (s PaymentStatus) String() string {
	return string(s)
}

func (s PaymentStatus) IsPaid() bool {
	return s == PaymentStatusPaid || s == PaymentStatusRefunded
}
func (s PaymentStatus) IsUnpaid() bool {
	return s == PaymentStatusUnpaid
}
func (s PaymentStatus) IsRefunded() bool {
	return s == PaymentStatusRefunded
}
