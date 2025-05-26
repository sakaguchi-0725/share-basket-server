package repository

import (
	"context"
	"sharebasket/domain/model"
)

type Category interface {
	GetAll(ctx context.Context) ([]model.Category, error)
}
