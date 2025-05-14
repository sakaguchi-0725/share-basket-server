package repository

import "sharebasket/domain/model"

type PersonalItem interface {
	GetAll(accID model.AccountID, status *model.ShoppingStatus) ([]model.PersonalItem, error)
	Store(item *model.PersonalItem) error
	GetByID(id int64) (model.PersonalItem, error)
	Delete(id int64) error
}
