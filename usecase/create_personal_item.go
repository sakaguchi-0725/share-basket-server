package usecase

import (
	"context"
	"sharebasket/core"
	"sharebasket/domain/model"
	"sharebasket/domain/repository"
)

type (
	CreatePersonalItem interface {
		Execute(ctx context.Context, in CreatePersonalItemInput) error
	}

	CreatePersonalItemInput struct {
		UserID     string
		Name       string
		Status     string
		CategoryID int64
	}

	createPersonalItem struct {
		accountRepo  repository.Account
		personalRepo repository.PersonalItem
	}
)

func (c *createPersonalItem) Execute(ctx context.Context, in CreatePersonalItemInput) error {
	account, err := c.accountRepo.Get(ctx, in.UserID)
	if err != nil {
		return err
	}

	var status model.ShoppingStatus
	if in.Status != "" {
		status, err = model.ParseShoppingStatus(in.Status)
		if err != nil {
			return core.NewInvalidError(err)
		}
	}

	item, err := model.NewPersonalItem(in.Name, core.Ptr(status), in.CategoryID, account.ID)
	if err != nil {
		return core.NewInvalidError(err)
	}

	if err := c.personalRepo.Store(&item); err != nil {
		return err
	}

	return nil
}

func NewCreatePersonalItem(a repository.Account, p repository.PersonalItem) CreatePersonalItem {
	return &createPersonalItem{
		accountRepo:  a,
		personalRepo: p,
	}
}
