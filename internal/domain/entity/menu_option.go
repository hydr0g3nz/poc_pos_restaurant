package entity

import (
	errs "github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/error"
)

// MenuOption represents a menu option domain entity
type MenuOption struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	IsRequired bool   `json:"isRequired"`
}

func (mo *MenuOption) IsValid() bool {
	return mo.Name != "" && mo.Type != ""
}

// NewMenuOption creates a new menu option
func NewMenuOption(name, optionType string, isRequired bool) (*MenuOption, error) {
	if name == "" || optionType == "" {
		return nil, errs.ErrInvalidMenuOption
	}

	return &MenuOption{
		Name:       name,
		Type:       optionType,
		IsRequired: isRequired,
	}, nil
}
