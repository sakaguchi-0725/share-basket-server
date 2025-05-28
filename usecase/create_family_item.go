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
	CreateFamilyItem interface {
		Execute(ctx context.Context, in CreateFamilyItemInput) error
	}

	CreateFamilyItemInput struct {
		UserID     string
		Name       string
		Status     string
		CategoryID int64
	}

	createFamilyItem struct {
		accountRepo    repository.Account
		familyRepo     repository.Family
		familyItemRepo repository.FamilyItem
		accountService service.Account
	}
)

func (c *createFamilyItem) Execute(ctx context.Context, in CreateFamilyItemInput) error {
	account, err := c.accountRepo.Get(ctx, in.UserID)
	if err != nil {
		return err
	}

	// 家族に参加しているか確認
	hasFamily, err := c.accountService.HasFamily(ctx, account.ID)
	if err != nil {
		return err
	}

	if !hasFamily {
		return core.NewInvalidError(errors.New("user is not a member of any family")).
			WithMessage("家族に参加していません")
	}

	family, err := c.familyRepo.GetByAccountID(ctx, account.ID)
	if err != nil {
		return err
	}

	status, err := in.ShoppingStatus()
	if err != nil {
		return err
	}

	item, err := model.NewFamilyItem(in.Name, status, in.CategoryID, family.ID, account.ID)
	if err != nil {
		return err
	}

	err = c.familyItemRepo.Store(ctx, &item)
	if err != nil {
		return err
	}

	return nil
}

func (in CreateFamilyItemInput) ShoppingStatus() (*model.ShoppingStatus, error) {
	if in.Status == "" {
		return nil, nil
	}

	status, err := model.ParseShoppingStatus(in.Status)
	if err != nil {
		return nil, err
	}

	return core.Ptr(status), nil
}

func NewCreateFamilyItem(
	a repository.Account,
	rf repository.Family,
	rfi repository.FamilyItem,
	as service.Account,
) CreateFamilyItem {
	return &createFamilyItem{
		accountRepo:    a,
		familyRepo:     rf,
		familyItemRepo: rfi,
		accountService: as,
	}
}
