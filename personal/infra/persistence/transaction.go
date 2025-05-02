package persistence

import (
	"context"
	"share-basket-server/personal/domain"

	"gorm.io/gorm"
)

type transaction struct {
	db *gorm.DB
}

func (t *transaction) Run(ctx context.Context, fn func(ctx context.Context) error) error {
	return t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(ctx)
	})
}

func NewTransaction(db *gorm.DB) domain.Transaction {
	return &transaction{db}
}
