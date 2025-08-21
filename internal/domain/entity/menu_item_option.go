package entity

// MenuItem represents a menu item domain entity
type MenuItemOption struct {
	ItemID   int  `json:"id"`
	OptionID int  `json:"optionId"`
	IsActive bool `json:"is_active"`
}

// IsValid checks whether the menu item option is valid
func (mio *MenuItemOption) IsValid() bool {
	return mio.ItemID > 0 && mio.OptionID > 0
}

func NewMenuItemOption(itemID, optionID int, isActive bool) (*MenuItemOption, error) {
	return &MenuItemOption{
		ItemID:   itemID,
		OptionID: optionID,
		IsActive: isActive,
	}, nil
}
