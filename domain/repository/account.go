package repository

import (
	"context"
	"sharebasket/domain/model"
)

type Account interface {
	Get(userID string) (model.Account, error)
	Store(ctx context.Context, acc *model.Account) error
}
