package entity

import (
	errs "github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/error"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/vo"
)

// MenuOption represents a menu option domain entity
type MenuOption struct {
	ID           int            `json:"id"`
	Name         string         `json:"name"`
	Type         vo.OptionType  `json:"type"`
	IsRequired   bool           `json:"isRequired"`
	OptionValues []*OptionValue `json:"optionValues"`
}

func (mo *MenuOption) IsValid() bool {
	return mo.Name != "" && mo.Type != ""
}

// NewMenuOption creates a new menu option
func NewMenuOption(name string, optionType vo.OptionType, isRequired bool) (*MenuOption, error) {
	if name == "" || optionType == "" {
		return nil, errs.ErrInvalidMenuOption
	}

	return &MenuOption{
		Name:       name,
		Type:       optionType,
		IsRequired: isRequired,
	}, nil
}
