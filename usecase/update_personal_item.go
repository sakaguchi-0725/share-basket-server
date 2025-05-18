package usecase

import (
	"context"
	"errors"
	"sharebasket/core"
	"sharebasket/domain/model"
	"sharebasket/domain/repository"
)

var ErrPersonalItemNotFound = errors.New("personal item is not found")

type (
	UpdatePersonalItem interface {
		Execute(ctx context.Context, in UpdatePersonalItemInput) error
	}

	UpdatePersonalItemInput struct {
		ID         int64
		Name       string
		Status     string
		CategoryID int64
		UserID     string
	}

	updatePersonalItem struct {
		accountRepo  repository.Account
		personalRepo repository.PersonalItem
		logger       core.Logger
	}
)

func (u *updatePersonalItem) Execute(ctx context.Context, in UpdatePersonalItemInput) error {
	account, err := u.accountRepo.Get(in.UserID)
	if err != nil {
		u.logger.WithError(err).
			With("user_id", in.UserID).
			Error("failed to get account")
		return err
	}

	item, err := u.personalRepo.GetByID(in.ID)
	if err != nil {
		if errors.Is(err, ErrPersonalItemNotFound) {
			u.logger.WithError(err).
				With("item_id", in.ID).
				Warn("personal item not found")
			return core.NewInvalidError(err)
		}

		u.logger.WithError(err).
			With("item_id", in.ID).
			Error("failed to get personal item")
		return err
	}

	// 買い物リストの所有権確認
	if err := item.CheckOwner(account.ID); err != nil {
		u.logger.WithError(err).
			With("account_id", account.ID.String()).
			With("item_id", in.ID).
			Warn("unauthorized delete attempt - item owner mismatch")
		return core.NewAppError(core.ErrForbidden, err)
	}

	// ステータスの更新
	var status *model.ShoppingStatus
	if in.Status != "" {
		s, err := model.ParseShoppingStatus(in.Status)
		if err != nil {
			u.logger.WithError(err).
				With("status", in.Status).
				Warn("invalid shopping status")
			return core.NewInvalidError(err)
		}
		status = &s
	}

	// カテゴリーIDの更新
	var categoryID *int64
	if in.CategoryID != 0 {
		categoryID = &in.CategoryID
	}

	// アイテムの更新
	if err := item.Update(in.Name, status, categoryID); err != nil {
		u.logger.WithError(err).
			With("name", in.Name).
			With("status", status.String()).
			With("category_id", categoryID).
			Warn("failed to update shopping item model")
		return core.NewInvalidError(err)
	}

	// 更新を保存
	if err := u.personalRepo.Store(&item); err != nil {
		u.logger.WithError(err).Error("failed to store shopping item")
		return err
	}

	return nil
}

func NewUpdatePersonalItem(a repository.Account, p repository.PersonalItem, l core.Logger) UpdatePersonalItem {
	return &updatePersonalItem{
		accountRepo:  a,
		personalRepo: p,
		logger:       l,
	}
}
