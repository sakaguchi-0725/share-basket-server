package domain

import (
	"errors"
	"share-basket-server/core/util"
)

var ErrShoppingItemNameRequired = errors.New("shopping item name is required")

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
	name string, status *ShoppingStatus, categoryID uint,
) (ShoppingItem, error) {
	if name == "" {
		return ShoppingItem{}, ErrShoppingItemNameRequired
	}

	var itemStatus ShoppingStatus
	if status == nil {
		itemStatus = UnPurchased
	} else {
		itemStatus = util.Derefer(status)
	}

	return ShoppingItem{
		Name:       name,
		Status:     itemStatus,
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
