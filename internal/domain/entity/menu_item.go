package entity

import (
	"time"

	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/vo"
)

// MenuItem represents a menu item domain entity
type MenuItem struct {
	ID              int       `json:"id"`
	CategoryID      int       `json:"category_id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Price           vo.Money  `json:"price"`
	ImageURL        string    `json:"image_url,omitempty"`
	IsRecommended   bool      `json:"is_recommended,omitempty"`
	DiscountPercent float64   `json:"discount_percent,omitempty"`
	IsDiscounted    bool      `json:"is_discounted,omitempty"`
	PreparationTime int       `json:"preparation_time,omitempty"` // in minutes
	DisplayOrder    int       `json:"display_order,omitempty"`
	KitchenID       int       `json:"kitchen_station_id,omitempty"` // optional kitchen station for tracking
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	// relationships
	Category        *Category         `json:"category,omitempty"`
	KitchenStation  *KitchenStation   `json:"kitchen_station,omitempty"`
	MenuItemOptions []*MenuItemOption `json:"menu_item_options,omitempty"`
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
