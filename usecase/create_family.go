package usecase

import (
	"context"
	"errors"
	"sharebasket/core"
	"sharebasket/domain/model"
	"sharebasket/domain/repository"
	"sharebasket/domain/service"
)

type (
	CreateFamily interface {
		Execute(ctx context.Context, in CreateFamilyInput) error
	}

	CreateFamilyInput struct {
		Name   string
		UserID string
	}

	createFamily struct {
		accountRepo    repository.Account
		familyRepo     repository.Family
		accountService service.Account
	}
)

func (c *createFamily) Execute(ctx context.Context, in CreateFamilyInput) error {
	account, err := c.accountRepo.Get(ctx, in.UserID)
	if err != nil {
		return err
	}

	// すでに家族オーナー、または家族メンバーではないか判定
	hasFamily, err := c.accountService.HasFamily(ctx, account.ID)
	if err != nil {
		return err
	}

	if hasFamily {
		return core.NewInvalidError(errors.New("account already has family"))
	}

	id := model.NewFamilyID()

	family, err := model.NewFamily(id, in.Name, account)
	if err != nil {
		return core.NewInvalidError(err)
	}

	err = c.familyRepo.Store(ctx, &family)
	if err != nil {
		return err
	}

	return nil
}

func NewCreateFamily(a repository.Account, f repository.Family, as service.Account) CreateFamily {
	return &createFamily{
		accountRepo:    a,
		familyRepo:     f,
		accountService: as,
	}
}
