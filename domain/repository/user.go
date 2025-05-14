package repository

import (
	"context"
	"sharebasket/domain/model"
)

type User interface {
	ExistsByEmail(email string) (bool, error)
	Store(ctx context.Context, user *model.User) error
}
