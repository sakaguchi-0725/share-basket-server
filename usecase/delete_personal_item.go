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
		logger       core.Logger
	}
)

func (d *deletePersonalItem) Execute(ctx context.Context, in DeletePersonalItemInput) error {
	account, err := d.accountRepo.Get(in.UserID)
	if err != nil {
		if errors.Is(err, ErrAccountNotFound) {
			d.logger.WithError(err).
				With("user_id", in.UserID).
				Warn("account not found")
			return core.NewAppError(core.ErrUnauthorized, err)
		}

		d.logger.WithError(err).
			With("user_id", in.UserID).
			Error("failed to get account")
		return err
	}

	item, err := d.personalRepo.GetByID(in.ID)
	if err != nil {
		if errors.Is(err, ErrPersonalItemNotFound) {
			d.logger.WithError(err).
				With("item_id", in.ID).
				Warn("personal item not found")
			return core.NewInvalidError(err)
		}

		d.logger.WithError(err).
			With("item_id", in.ID).
			Error("failed to get personal item")
		return err
	}

	// 買い物リストの所有権確認
	if err := item.CheckOwner(account.ID); err != nil {
		d.logger.WithError(err).
			With("account_id", account.ID.String()).
			With("item_id", in.ID).
			Warn("unauthorized delete attempt - item owner mismatch")
		return core.NewAppError(core.ErrForbidden, err)
	}

	if err := d.personalRepo.Delete(*item.ID); err != nil {
		d.logger.WithError(err).
			With("item_id", in.ID).
			Error("failed to delete personal item from repository")
		return err
	}

	return nil
}

func NewDeletePersonalItem(a repository.Account, p repository.PersonalItem, l core.Logger) DeletePersonalItem {
	return &deletePersonalItem{
		accountRepo:  a,
		personalRepo: p,
		logger:       l,
	}
}
