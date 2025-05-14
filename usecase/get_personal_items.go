package usecase

import (
	"context"
	"errors"
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
		logger       core.Logger
	}
)

func (g *getPersonalItems) Execute(ctx context.Context, in GetPersonalItemsInput) ([]GetPersonalItemsOutput, error) {
	account, err := g.accountRepo.Get(in.UserID)
	if err != nil {
		if errors.Is(err, ErrAccountNotFound) {
			g.logger.WithError(err).
				With("user_id", in.UserID).
				Warn("account not found")
			return []GetPersonalItemsOutput{}, core.NewAppError(core.ErrUnauthorized, err)
		}

		g.logger.WithError(err).
			With("user_id", in.UserID).
			Error("failed to get account")
		return []GetPersonalItemsOutput{}, err
	}

	status, err := in.status()
	if err != nil {
		g.logger.WithError(err).
			With("status", in.Status).
			Warn("invalid shopping status")
		return []GetPersonalItemsOutput{}, err
	}

	items, err := g.personalRepo.GetAll(account.ID, status)
	if err != nil {
		g.logger.WithError(err).
			Error("failed to get personal shopping item")
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

func NewGetPersonalItems(a repository.Account, p repository.PersonalItem, l core.Logger) GetPersonalItems {
	return &getPersonalItems{
		accountRepo:  a,
		personalRepo: p,
		logger:       l,
	}
}
