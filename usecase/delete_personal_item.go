package usecase

import (
	"context"
	"errors"
	"sharebasket/core"
	"sharebasket/domain/repository"
)

type (
	DeletePersonalItem interface {
		Execute(ctx context.Context, in DeletePersonalItemInput) error
	}

	DeletePersonalItemInput struct {
		ID     int64
		UserID string
	}

	deletePersonalItem struct {
		accountRepo  repository.Account
		personalRepo repository.PersonalItem
	}
)

func (d *deletePersonalItem) Execute(ctx context.Context, in DeletePersonalItemInput) error {
	account, err := d.accountRepo.Get(ctx, in.UserID)
	if err != nil {
		return err
	}

	item, err := d.personalRepo.GetByID(in.ID)
	if err != nil {
		if errors.Is(err, ErrPersonalItemNotFound) {
			return core.NewInvalidError(err)
		}

		return err
	}

	// 買い物リストの所有権確認
	if err := item.CheckOwner(account.ID); err != nil {
		return core.NewAppError(core.ErrForbidden, err)
	}

	if err := d.personalRepo.Delete(*item.ID); err != nil {
		return err
	}

	return nil
}

func NewDeletePersonalItem(a repository.Account, p repository.PersonalItem) DeletePersonalItem {
	return &deletePersonalItem{
		accountRepo:  a,
		personalRepo: p,
	}
}
