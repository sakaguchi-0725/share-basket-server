package usecase

import (
	"context"
	"sharebasket/core"
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
		repo   repository.Category
		logger core.Logger
	}
)

func (g *getCategories) Execute(ctx context.Context) ([]GetCategoriesOutput, error) {
	categories, err := g.repo.GetAll()
	if err != nil {
		g.logger.WithError(err).
			Error("failed to get category")
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

func NewGetCategories(r repository.Category, l core.Logger) GetCategories {
	return &getCategories{
		repo:   r,
		logger: l,
	}
}
