package vo

import (
	"strings"

	errs "github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/error"
)

type PaymentMethod string

const (
	PaymentMethodCash       PaymentMethod = "cash"
	PaymentMethodCreditCard PaymentMethod = "credit_card"
	PaymentMethodWallet     PaymentMethod = "wallet"
)

func (p PaymentMethod) Valid() bool {
	switch p {
	case PaymentMethodCash, PaymentMethodCreditCard, PaymentMethodWallet:
		return true
	default:
		return false
	}
}

func NewPaymentMethod(method string) (PaymentMethod, error) {
	pm := PaymentMethod(strings.ToLower(method))
	if !pm.Valid() {
		return "", errs.ErrInvalidPaymentMethod
	}
	return pm, nil
}

func (p PaymentMethod) String() string {
	return string(p)
}
