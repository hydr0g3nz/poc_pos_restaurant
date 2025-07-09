package vo

import (
	"fmt"

	errs "github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/error"
)

type Price struct {
	amount float64
}

func NewPrice(amount float64) (Price, error) {
	if amount < 0 {
		return Price{}, errs.ErrInvalidMenuItemPrice
	}
	if amount > 100000 { // Maximum price 100,000
		return Price{}, errs.ErrInvalidMenuItemPrice
	}

	return Price{amount: amount}, nil
}

func (p Price) Amount() float64 {
	return p.amount
}

func (p Price) String() string {
	return fmt.Sprintf("%.2f", p.amount)
}

func (p Price) ToMoney() Money {
	money, _ := NewMoney(p.amount)
	return money
}
