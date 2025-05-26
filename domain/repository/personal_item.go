package repository

import (
	"context"
	"sharebasket/domain/model"
)

type PersonalItem interface {
	GetAll(ctx context.Context, accID model.AccountID, status *model.ShoppingStatus) ([]model.PersonalItem, error)
	Store(ctx context.Context, item *model.PersonalItem) error
	GetByID(ctx context.Context, id int64) (model.PersonalItem, error)
	Delete(ctx context.Context, id int64) error
}
