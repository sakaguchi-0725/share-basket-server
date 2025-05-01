package repository

import "share-basket-server/personal/domain/model"

type ShoppingCategory interface {
	GetAll() ([]model.ShoppingCategory, error)
}
