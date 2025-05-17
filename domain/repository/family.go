package repository

import "sharebasket/domain/model"

type Family interface {
	Store(family *model.Family) error
	Join(family *model.Family) error
	HasFamily(accountID model.AccountID) (bool, error)
}
