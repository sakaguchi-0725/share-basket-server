package repository

import (
	"context"
	"sharebasket/domain/model"
)

type FamilyItem interface {
	Get(ctx context.Context, id model.FamilyID, status *model.ShoppingStatus) ([]model.FamilyItem, error)
	Store(ctx context.Context, item *model.FamilyItem) error
}
