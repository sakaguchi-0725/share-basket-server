package repository

import "share-basket/personal-shopping/domain/model"

type ShoppingItemRepository interface {
	Get(status model.ShoppingStatus) ([]model.ShoppingItem, error)
	Save(item *model.ShoppingItem) error
}
