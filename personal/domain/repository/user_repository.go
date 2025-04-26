package repository

import "share-basket-server/personal/domain/model"

type User interface {
	GetByEmail(email string) (model.User, error)
	Store(user *model.User) error
}
