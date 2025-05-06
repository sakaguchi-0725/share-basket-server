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
		AccountID  AccountID
	}

	PersonalShoppingItemRepository interface {
		GetAll(accID AccountID, status *ShoppingStatus) ([]PersonalShoppingItem, error)
		Store(item *PersonalShoppingItem) error
	}
)

func NewPersonalShoppingItem(
	name string, status *ShoppingStatus, categoryID uint, accID AccountID,
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
		AccountID:  accID,
	}, nil
}
