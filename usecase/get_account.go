package usecase

import (
	"context"
	"sharebasket/domain/model"
	"sharebasket/domain/repository"
)

type (
	GetAccount interface {
		Execute(ctx context.Context, userID string) (GetAccountOutput, error)
	}

	GetAccountOutput struct {
		ID   string
		Name string
	}

	getAccount struct {
		repo repository.Account
	}
)

func (g *getAccount) Execute(ctx context.Context, userID string) (GetAccountOutput, error) {
	account, err := g.repo.Get(ctx, userID)
	if err != nil {
		return GetAccountOutput{}, err
	}

	return g.makeOutput(account), nil
}

func (g *getAccount) makeOutput(acc model.Account) GetAccountOutput {
	return GetAccountOutput{
		ID:   acc.ID.String(),
		Name: acc.Name,
	}
}

func NewGetAccount(r repository.Account) GetAccount {
	return &getAccount{
		repo: r,
	}
}
