package entity

import (
	errs "github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/error"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/vo"
)

// OptionValue represents a menu option value
type OptionValue struct {
	ID              int      `json:"id"`
	OptionID        int      `json:"optionId"`
	Name            string   `json:"name"`
	IsDefault       bool     `json:"isDefault"`
	AdditionalPrice vo.Money `json:"additionalPrice,omitempty"` // optional additional price for this option value
	DisplayOrder    int      `json:"displayOrder,omitempty"`    // optional display order for sorting
}

// IsValid checks whether the option value is valid
func (ov *OptionValue) IsValid() bool {
	return ov.OptionID > 0 && ov.Name != ""
}

// NewOptionValue creates a new option value
func NewOptionValue(optionID int, name string, isDefault bool) (*OptionValue, error) {
	if optionID <= 0 || name == "" {
		return nil, errs.ErrInvalidMenuOptionValue
	}

	return &OptionValue{
		OptionID:  optionID,
		Name:      name,
		IsDefault: isDefault,
	}, nil
}
