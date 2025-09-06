package entity

// Table represents a restaurant table domain entity
type Table struct {
	ID           int                `json:"id"`
	TableNumber  int                `json:"table_number"`
	Seating      int                `json:"seating"`
	IsActive     bool               `json:"is_active"`
	CurrentOrder *OrderTableDetails `json:"order_details,omitempty"`
}
type OrderTableDetails struct {
	OrderID     int    `json:"order_id"`
	OrderNumber int    `json:"order_number"`
	Status      string `json:"status"`
	QRCode      string `json:"qr_code"`
	CreatedAt   string `json:"created_at"`
}

// IsValid validates table data
func (t *Table) IsValid() bool {
	return t.TableNumber > 0 && t.Seating >= 0
}

// NewTable creates a new table
func NewTable(tableNumber int, seating int) (*Table, error) {

	return &Table{
		TableNumber: tableNumber,
		Seating:     seating,
		IsActive:    true,
	}, nil
}
