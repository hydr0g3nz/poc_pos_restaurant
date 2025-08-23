package entity

type KitchenStation struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	IsAvailable bool   `json:"isAvailable"`
}

func NewKitchenStation(name string, isAvailable bool) (*KitchenStation, error) {
	return &KitchenStation{
		Name:        name,
		IsAvailable: isAvailable,
	}, nil
}
