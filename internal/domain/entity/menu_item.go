package entity

import (
	"time"

	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/vo"
)

// MenuItem represents a menu item domain entity
type MenuItem struct {
	ID          int       `json:"id"`
	CategoryID  int       `json:"category_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       vo.Money  `json:"price"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// IsValid validates menu item data
func (m *MenuItem) IsValid() bool {
	return m.Name != "" && m.CategoryID > 0 && !m.Price.IsZero()
}

// NewMenuItem creates a new menu item
func NewMenuItem(categoryID int, name, description string, price float64) (*MenuItem, error) {
	money, err := vo.NewMoneyFromBaht(price)
	if err != nil {
		return nil, err
	}

	return &MenuItem{
		CategoryID:  categoryID,
		Name:        name,
		Description: description,
		Price:       money,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

// Activate activates the menu item
func (m *MenuItem) Activate() {
	m.IsActive = true
	m.UpdatedAt = time.Now()
}

// Deactivate deactivates the menu item
func (m *MenuItem) Deactivate() {
	m.IsActive = false
	m.UpdatedAt = time.Now()
}

// UpdatePrice updates the price of the menu item
func (m *MenuItem) UpdatePrice(newPrice float64) error {
	money, err := vo.NewMoneyFromBaht(newPrice)
	if err != nil {
		return err
	}

	m.Price = money
	m.UpdatedAt = time.Now()
	return nil
}
