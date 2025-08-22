package entity

import (
	"time"

	errs "github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/error"
)

// Category represents a menu category domain entity
type Category struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description,omitempty"`
	DisplayOrder int       `json:"display_order,omitempty"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// IsValid validates category data
func (c *Category) IsValid() bool {
	return c.Name != ""
}

// NewCategory creates a new category
func NewCategory(name string, description string, displayOrder int, isActive bool) (*Category, error) {
	if name == "" {
		return nil, errs.ErrInvalidCategoryName
	}
	return &Category{
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsActive:  isActive,
	}, nil
}
