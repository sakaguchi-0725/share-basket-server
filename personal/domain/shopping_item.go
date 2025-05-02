package domain

import "errors"

type (
	ShoppingItem struct {
		ID         *uint
		Name       string
		Status     ShoppingStatus
		CategoryID uint
	}

	ShoppingItemRepository interface {
		GetAll() ([]ShoppingItem, error)
		Store(item *ShoppingItem) error
	}
)

func NewShoppingItem(
	id *uint, name string, status ShoppingStatus, categoryID uint,
) (ShoppingItem, error) {
	if id == nil {
		return ShoppingItem{}, errors.New("shopping id is required")
	}

	if name == "" {
		return ShoppingItem{}, errors.New("name is required")
	}

	return ShoppingItem{
		ID:         id,
		Name:       name,
		Status:     status,
		CategoryID: categoryID,
	}, nil
}

func CreateShoppingItem(name string, categoryID uint) (ShoppingItem, error) {
	if name == "" {
		return ShoppingItem{}, errors.New("name is required")
	}

	return ShoppingItem{
		Name:       name,
		Status:     UnPurchased,
		CategoryID: categoryID,
	}, nil
}

func RecreateShoppingItem(id *uint, name string, status ShoppingStatus, categoryID uint) ShoppingItem {
	return ShoppingItem{
		ID:         id,
		Name:       name,
		Status:     status,
		CategoryID: categoryID,
	}
}
