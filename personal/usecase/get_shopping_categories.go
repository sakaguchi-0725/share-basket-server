package usecase

import (
	"context"
	"share-basket-server/personal/domain/repository"
	"share-basket-server/personal/usecase/input"
	"share-basket-server/personal/usecase/output"
)

type getShoppingCategories struct {
	repo repository.ShoppingCategory
}

func (usecase *getShoppingCategories) Execute(ctx context.Context, out output.GetShoppingCategoriesPort) error {
	categories, err := usecase.repo.GetAll()
	if err != nil {
		return err
	}

	return out.Render(ctx, output.MakeShoppingCategories(categories))
}

func NewGetShoppingCategories(repo repository.ShoppingCategory) input.GetShoppingCategoriesPort {
	return &getShoppingCategories{repo}
}
