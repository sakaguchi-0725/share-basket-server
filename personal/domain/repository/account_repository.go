package repository

import "share-basket-server/personal/domain/model"

type Account interface {
	Store(acc *model.Account) error
}
