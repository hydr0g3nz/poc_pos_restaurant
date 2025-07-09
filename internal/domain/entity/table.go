package entity

import (
	"fmt"

	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/vo"
)

// Table represents a restaurant table domain entity
type Table struct {
	ID          int            `json:"id"`
	TableNumber vo.TableNumber `json:"table_number"`
	QRCode      string         `json:"qr_code"`
	Seating     int            `json:"seating"`
	IsActive    bool           `json:"is_active"`
}

// IsValid validates table data
func (t *Table) IsValid() bool {
	return t.TableNumber.Number() > 0 && t.QRCode != "" && t.Seating >= 0
}

// NewTable creates a new table
func NewTable(tableNumber int, seating int) (*Table, error) {
	tableNum, err := vo.NewTableNumber(tableNumber)
	if err != nil {
		return nil, err
	}

	return &Table{
		TableNumber: tableNum,
		QRCode:      fmt.Sprintf("/order?table=%d", tableNumber),
		Seating:     seating,
		IsActive:    true,
	}, nil
}

// GenerateQRCode generates QR code URL for table
func (t *Table) GenerateQRCode() string {
	return fmt.Sprintf("/order?table=%d", t.TableNumber.Number())
}

// UpdateQRCode updates the QR code
func (t *Table) UpdateQRCode() {
	t.QRCode = t.GenerateQRCode()
}
