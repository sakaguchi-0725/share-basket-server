package usecase

import (
	"context"
	"sharebasket/domain/model"
	"sharebasket/domain/repository"
)

type (
	GetCategories interface {
		Execute(ctx context.Context) ([]GetCategoriesOutput, error)
	}

	GetCategoriesOutput struct {
		ID   int64
		Name string
	}

	getCategories struct {
		repo repository.Category
	}
)

func (g *getCategories) Execute(ctx context.Context) ([]GetCategoriesOutput, error) {
	categories, err := g.repo.GetAll(ctx)
	if err != nil {
		return []GetCategoriesOutput{}, err
	}

	return g.makeOutput(categories), nil
}

func (g *getCategories) makeOutput(categories []model.Category) []GetCategoriesOutput {
	outputs := make([]GetCategoriesOutput, len(categories))

	for i, v := range categories {
		outputs[i] = GetCategoriesOutput{
			ID:   v.ID,
			Name: v.Name,
		}
	}

	return outputs
}

func NewGetCategories(r repository.Category) GetCategories {
	return &getCategories{
		repo: r,
	}
}
