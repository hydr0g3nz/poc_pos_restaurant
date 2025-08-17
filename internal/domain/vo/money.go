package vo

import (
	"encoding/json"
	"fmt"

	err "github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/error"
)

type Money struct {
	// เก็บหน่วยเป็น "สตางค์"
	amount int64
}

func NewMoneyFromBaht(amount float64) (Money, error) {
	if amount < 0 {
		return Money{}, err.ErrNegativeAmount
	}
	return Money{amount: int64(amount*100 + 0.5)}, nil // ปัดเศษ
}

func NewMoneyFromSatang(amount int64) (Money, error) {
	if amount < 0 {
		return Money{}, err.ErrNegativeAmount
	}
	return Money{amount: amount}, nil
}

func (m Money) AmountSatang() int64 {
	return m.amount
}

func (m Money) AmountBaht() float64 {
	return float64(m.amount) / 100.0
}

func (m Money) Add(other Money) Money {
	return Money{amount: m.amount + other.amount}
}

func (m Money) Subtract(other Money) (Money, error) {
	if m.amount < other.amount {
		return Money{}, err.ErrInsufficientBalance
	}
	return Money{amount: m.amount - other.amount}, nil
}

func (m Money) Multiply(n float64) Money {
	return Money{amount: int64(float64(m.amount) * n)}
}

func (m Money) Divide(n float64) (Money, error) {
	if n == 0 {
		return Money{}, fmt.Errorf("division by zero")
	}
	return Money{amount: int64(float64(m.amount) / n)}, nil
}

func (m Money) IsZero() bool {
	return m.amount == 0
}

func (m Money) String() string {
	baht := m.amount / 100
	satang := m.amount % 100
	return fmt.Sprintf("%d.%02d", baht, satang)
}

// JSON serialize เป็น string เช่น "100.25"
func (m Money) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.String())
}

func (m *Money) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	var baht int64
	var satang int64
	if _, err := fmt.Sscanf(s, "%d.%02d", &baht, &satang); err != nil {
		return err
	}
	m.amount = baht*100 + satang
	return nil
}
