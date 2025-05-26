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
	}
)

func (u *updatePersonalItem) Execute(ctx context.Context, in UpdatePersonalItemInput) error {
	account, err := u.accountRepo.Get(ctx, in.UserID)
	if err != nil {
		return err
	}

	item, err := u.personalRepo.GetByID(in.ID)
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

	// ステータスの更新
	var status *model.ShoppingStatus
	if in.Status != "" {
		s, err := model.ParseShoppingStatus(in.Status)
		if err != nil {
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
		return core.NewInvalidError(err)
	}

	// 更新を保存
	if err := u.personalRepo.Store(&item); err != nil {
		return err
	}

	return nil
}

func NewUpdatePersonalItem(a repository.Account, p repository.PersonalItem) UpdatePersonalItem {
	return &updatePersonalItem{
		accountRepo:  a,
		personalRepo: p,
	}
}
