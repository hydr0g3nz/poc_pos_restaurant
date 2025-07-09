package vo

import (
	"fmt"

	errs "github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/error"
)

type Quantity struct {
	value int
}

func NewQuantity(value int) (Quantity, error) {
	if value <= 0 {
		return Quantity{}, errs.ErrInvalidQuantity
	}
	return Quantity{value: value}, nil
}

func (q Quantity) Value() int {
	return q.value
}

func (q Quantity) Add(other Quantity) Quantity {
	return Quantity{value: q.value + other.value}
}

func (q Quantity) Subtract(other Quantity) (Quantity, error) {
	if q.value < other.value {
		return Quantity{}, errs.ErrInsufficientQuantity
	}
	return Quantity{value: q.value - other.value}, nil
}

func (q Quantity) String() string {
	return fmt.Sprintf("%d", q.value)
}
