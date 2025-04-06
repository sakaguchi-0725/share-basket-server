package model

import "errors"

type ShoppingItem struct {
	ID       *int64
	Name     string
	Category ShoppingCategory
	Status   ShoppingStatus
}

func NewShoppingItem(name string, category ShoppingCategory, status ShoppingStatus) (ShoppingItem, error) {
	if name == "" {
		return ShoppingItem{}, errors.New("shopping name is required")
	}

	return ShoppingItem{
		ID:       nil,
		Name:     name,
		Category: category,
		Status:   status,
	}, nil
}

func RecreateShoppingItem(id *int64, name string, category ShoppingCategory, status ShoppingStatus) ShoppingItem {
	return ShoppingItem{
		ID:       id,
		Name:     name,
		Category: category,
		Status:   status,
	}
}
