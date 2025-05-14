package usecase

import (
	"context"
	"errors"
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
		logger       core.Logger
	}
)

// ユーザーの個人アイテムを作成するメソッドです。
// 指定されたユーザーIDに対応するアカウントを検証し、新しい個人アイテムを作成してリポジトリに保存します。
// アカウントが見つからない場合は認証エラーを、入力が無効な場合は不正エラーを返します。
// 成功した場合はnilを返します。
func (c *createPersonalItem) Execute(ctx context.Context, in CreatePersonalItemInput) error {
	account, err := c.accountRepo.Get(in.UserID)
	if err != nil {
		if errors.Is(err, ErrAccountNotFound) {
			c.logger.WithError(err).
				With("user_id", in.UserID).
				Warn("account not found")
			return core.NewAppError(core.ErrUnauthorized, err)
		}

		c.logger.WithError(err).
			With("user_id", in.UserID).
			Error("failed to get account")
		return err
	}

	var status model.ShoppingStatus
	if in.Status != "" {
		status, err = model.ParseShoppingStatus(in.Status)
		if err != nil {
			c.logger.WithError(err).
				With("status", in.Status).
				Warn("invalid shopping status")
			return core.NewInvalidError(err)
		}
	}

	item, err := model.NewPersonalItem(in.Name, core.Ptr(status), in.CategoryID, account.ID)
	if err != nil {
		c.logger.WithError(err).
			With("name", in.Name).
			With("category_id", in.CategoryID).
			Warn("invalid personal item parameters")
		return core.NewInvalidError(err)
	}

	if err := c.personalRepo.Store(&item); err != nil {
		c.logger.WithError(err).
			With("item", item).
			Error("failed to store personal item")
		return err
	}

	return nil
}

func NewCreatePersonalItem(a repository.Account, p repository.PersonalItem, l core.Logger) CreatePersonalItem {
	return &createPersonalItem{
		accountRepo:  a,
		personalRepo: p,
		logger:       l,
	}
}
