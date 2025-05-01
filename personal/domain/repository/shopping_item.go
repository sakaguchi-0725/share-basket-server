package repository

import "share-basket-server/personal/domain/model"

type ShoppingItem interface {
	GetAll() ([]model.ShoppingItem, error)
	Store(item *model.ShoppingItem) error
}
