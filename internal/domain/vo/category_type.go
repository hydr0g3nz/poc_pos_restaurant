package vo

import (
	"strings"

	errs "github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/error"
)

type CategoryType string

const (
	CategorySavory  CategoryType = "ของคาว"
	CategoryDessert CategoryType = "ของหวาน"
	CategorySnack   CategoryType = "ของทานเล่น"
	CategoryRoti    CategoryType = "โรตี"
)

func (c CategoryType) IsValid() bool {
	switch c {
	case CategorySavory, CategoryDessert, CategorySnack, CategoryRoti:
		return true
	default:
		return false
	}
}

func NewCategoryType(category string) (CategoryType, error) {
	c := CategoryType(strings.TrimSpace(category))
	if !c.IsValid() {
		return "", errs.ErrInvalidCategoryType
	}
	return c, nil
}

func (c CategoryType) String() string {
	return string(c)
}

// GetAllCategoryTypes returns all valid category types
func GetAllCategoryTypes() []CategoryType {
	return []CategoryType{
		CategorySavory,
		CategoryDessert,
		CategorySnack,
		CategoryRoti,
	}
}
