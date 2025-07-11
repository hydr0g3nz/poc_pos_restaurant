package entity

import (
	"time"

	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/vo"
)

// Category represents a menu category domain entity
type Category struct {
	ID        int             `json:"id"`
	Name      vo.CategoryType `json:"name"`
	CreatedAt time.Time       `json:"created_at"`
}

// IsValid validates category data
func (c *Category) IsValid() bool {
	return c.Name.IsValid()
}

// NewCategory creates a new category
func NewCategory(name string) (*Category, error) {
	categoryType, err := vo.NewCategoryType(name)
	if err != nil {
		return nil, err
	}

	return &Category{
		Name:      categoryType,
		CreatedAt: time.Now(),
	}, nil
}
