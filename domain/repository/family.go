package repository

import (
	"context"
	"sharebasket/domain/model"
)

type Family interface {
	Store(ctx context.Context, family *model.Family) error
	GetByToken(ctx context.Context, token string) (model.Family, error)
	GetByAccountID(ctx context.Context, id model.AccountID) (model.Family, error)
	HasFamily(ctx context.Context, accountID model.AccountID) (bool, error)
	Invitation(ctx context.Context, token string, familyID model.FamilyID) error
	HasOwnedFamily(ctx context.Context, accountID model.AccountID) (bool, error)
	GetOwnedFamily(ctx context.Context, accountID model.AccountID) (model.Family, error)
	HasMembership(ctx context.Context, accountID model.AccountID, familyID model.FamilyID) (bool, error)
}
