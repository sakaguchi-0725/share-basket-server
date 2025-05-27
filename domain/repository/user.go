package repository

import (
	"context"
	"sharebasket/domain/model"
)

type User interface {
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	Store(ctx context.Context, user *model.User) error
}
