package usecase

import (
	"context"
	"sharebasket/core"
	"sharebasket/domain/model"
	"sharebasket/domain/repository"
)

type (
	GetPersonalItems interface {
		Execute(ctx context.Context, in GetPersonalItemsInput) ([]GetPersonalItemsOutput, error)
	}

	GetPersonalItemsInput struct {
		UserID string
		Status string
	}

	GetPersonalItemsOutput struct {
		ID         int64
		Name       string
		Status     string
		CategoryID int64
	}

	getPersonalItems struct {
		accountRepo  repository.Account
		personalRepo repository.PersonalItem
	}
)

func (g *getPersonalItems) Execute(ctx context.Context, in GetPersonalItemsInput) ([]GetPersonalItemsOutput, error) {
	account, err := g.accountRepo.Get(in.UserID)
	if err != nil {
		return []GetPersonalItemsOutput{}, err
	}

	status, err := in.status()
	if err != nil {
		return []GetPersonalItemsOutput{}, err
	}

	items, err := g.personalRepo.GetAll(account.ID, status)
	if err != nil {
		return []GetPersonalItemsOutput{}, err
	}

	return makeGetPersonalItemsOutput(items), nil
}

func (in GetPersonalItemsInput) status() (*model.ShoppingStatus, error) {
	if in.Status == "" {
		return nil, nil
	}

	status, err := model.ParseShoppingStatus(in.Status)
	if err != nil {
		return nil, core.NewInvalidError(err)
	}

	return &status, nil
}

func makeGetPersonalItemsOutput(items []model.PersonalItem) []GetPersonalItemsOutput {
	outputs := make([]GetPersonalItemsOutput, len(items))

	for i, v := range items {
		outputs[i] = GetPersonalItemsOutput{
			ID:         core.Derefer(v.ID),
			Name:       v.Name,
			Status:     v.Status.String(),
			CategoryID: v.CategoryID,
		}
	}

	return outputs
}

func NewGetPersonalItems(a repository.Account, p repository.PersonalItem) GetPersonalItems {
	return &getPersonalItems{
		accountRepo:  a,
		personalRepo: p,
	}
}
