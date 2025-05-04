//go:generate mockgen -destination=../mock/$GOPACKAGE/$GOFILE . PersonalShoppingItemRepository
package domain

import (
	"errors"
	"share-basket-server/core/util"
)

var ErrPersonalShoppingItemNameRequired = errors.New("shopping item name is required")

type (
	PersonalShoppingItem struct {
		ID         *uint
		Name       string
		Status     ShoppingStatus
		CategoryID uint
	}

	PersonalShoppingItemRepository interface {
		GetAll() ([]PersonalShoppingItem, error)
		Store(item *PersonalShoppingItem) error
	}
)

func NewPersonalShoppingItem(
	name string, status *ShoppingStatus, categoryID uint,
) (PersonalShoppingItem, error) {
	if name == "" {
		return PersonalShoppingItem{}, ErrPersonalShoppingItemNameRequired
	}

	var itemStatus ShoppingStatus
	if status == nil {
		itemStatus = UnPurchased
	} else {
		itemStatus = util.Derefer(status)
	}

	return PersonalShoppingItem{
		Name:       name,
		Status:     itemStatus,
		CategoryID: categoryID,
	}, nil
}

func RecreatePersonalShoppingItem(id *uint, name string, status ShoppingStatus, categoryID uint) PersonalShoppingItem {
	return PersonalShoppingItem{
		ID:         id,
		Name:       name,
		Status:     status,
		CategoryID: categoryID,
	}
}
