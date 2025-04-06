package repository

import "share-basket/personal-shopping/domain/model"

type ShoppingCategoryRepository interface {
	Get() ([]model.ShoppingCategory, error)
	GetByID(id int64) (model.ShoppingCategory, error)
}
