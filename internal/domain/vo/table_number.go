package vo

import (
	"fmt"

	errs "github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/error"
)

type TableNumber struct {
	number int
}

func NewTableNumber(number int) (TableNumber, error) {
	if number <= 0 {
		return TableNumber{}, errs.ErrInvalidTableNumber
	}
	return TableNumber{number: number}, nil
}

func (t TableNumber) Number() int {
	return t.number
}

func (t TableNumber) String() string {
	return fmt.Sprintf("Table %d", t.number)
}
