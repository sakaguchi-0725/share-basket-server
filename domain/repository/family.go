package repository

import (
	"context"
	"sharebasket/domain/model"
)

type Family interface {
	Store(family *model.Family) error
	Join(family *model.Family) error
	HasFamily(accountID model.AccountID) (bool, error)
	Invitation(ctx context.Context, token string, familyID model.FamilyID) error
	HasOwnedFamily(accountID model.AccountID) (bool, error)
	GetOwnedFamily(accountID model.AccountID) (model.Family, error)
}
