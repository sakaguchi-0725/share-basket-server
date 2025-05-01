package output

import (
	"context"
	"share-basket-server/personal/domain/model"
)

type ShoppingCategory struct {
	ID   uint
	Name string
}

type GetShoppingCategories []ShoppingCategory

type GetShoppingCategoriesPort interface {
	Render(ctx context.Context, out GetShoppingCategories) error
}

func MakeShoppingCategories(models []model.ShoppingCategory) GetShoppingCategories {
	categories := make(GetShoppingCategories, len(models))

	for i, v := range models {
		categories[i] = ShoppingCategory{
			ID:   v.ID,
			Name: v.Name,
		}
	}

	return categories
}
