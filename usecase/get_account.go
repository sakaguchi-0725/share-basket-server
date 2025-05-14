package usecase

import (
	"context"
	"errors"
	"sharebasket/core"
	"sharebasket/domain/model"
	"sharebasket/domain/repository"
)

var ErrAccountNotFound = errors.New("account not found")

type (
	GetAccount interface {
		Execute(ctx context.Context, userID string) (GetAccountOutput, error)
	}

	GetAccountOutput struct {
		ID   string
		Name string
	}

	getAccount struct {
		repo   repository.Account
		logger core.Logger
	}
)

func (g *getAccount) Execute(ctx context.Context, userID string) (GetAccountOutput, error) {
	account, err := g.repo.Get(userID)
	if err != nil {
		if errors.Is(err, ErrAccountNotFound) {
			g.logger.WithError(err).
				With("user_id", userID).
				Warn("account not found")
			return GetAccountOutput{}, core.NewAppError(core.ErrUnauthorized, err)
		}

		g.logger.WithError(err).
			Error("failed to get account")
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

func NewGetAccount(r repository.Account, l core.Logger) GetAccount {
	return &getAccount{
		repo:   r,
		logger: l,
	}
}
