package vo

import (
	errs "github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/error"
)

type ItemStatus string

const (
	ItemStatusPending   ItemStatus = "pending"
	ItemStatusPreparing ItemStatus = "preparing"
	ItemStatusReady     ItemStatus = "ready"
	ItemStatusServed    ItemStatus = "served"
	ItemStatusCancelled ItemStatus = "cancelled"
)

func (s ItemStatus) IsValid() bool {
	switch s {
	case ItemStatusPending, ItemStatusPreparing, ItemStatusReady, ItemStatusServed, ItemStatusCancelled:
		return true
	default:
		return false
	}
}

func (s ItemStatus) String() string {
	return string(s)
}

func NewItemStatus(status string) (ItemStatus, error) {
	switch status {
	case string(ItemStatusPending), string(ItemStatusPreparing), string(ItemStatusReady), string(ItemStatusServed), string(ItemStatusCancelled):
		return ItemStatus(status), nil
	default:
		return "", errs.ErrInvalidItemStatus
	}
}
