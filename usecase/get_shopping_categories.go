//go:generate mockgen -destination=../test/mock/usecase/get_shopping_categories_input.go . GetShoppingCategoriesInputPort
//go:generate mockgen -destination=../test/mock/usecase/get_shopping_categories_output.go . GetShoppingCategoriesOutputPort
package usecase

import (
	"context"
	"share-basket-server/domain"
)

type (
	GetShoppingCategoriesInputPort interface {
		Execute(ctx context.Context, output GetShoppingCategoriesOutputPort) error
	}

	GetShoppingCategoriesOutputPort interface {
		Render(ctx context.Context, outputs []GetShoppingCategoryOutput) error
	}

	GetShoppingCategoryOutput struct {
		ID   uint
		Name string
	}

	getShoppingCategoriesInteractor struct {
		repo domain.ShoppingCategoryRepository
	}
)

func (usecase *getShoppingCategoriesInteractor) Execute(ctx context.Context, output GetShoppingCategoriesOutputPort) error {
	categories, err := usecase.repo.GetAll()
	if err != nil {
		return err
	}

	return output.Render(ctx, usecase.makeOutputs(categories))
}

func (usecase *getShoppingCategoriesInteractor) makeOutputs(models []domain.ShoppingCategory) []GetShoppingCategoryOutput {
	categories := make([]GetShoppingCategoryOutput, len(models))

	for i, v := range models {
		categories[i] = GetShoppingCategoryOutput{
			ID:   v.ID,
			Name: v.Name,
		}
	}

	return categories
}

func NewGetShoppingCategoriesInteractor(repo domain.ShoppingCategoryRepository) GetShoppingCategoriesInputPort {
	return &getShoppingCategoriesInteractor{repo}
}
