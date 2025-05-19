package usecase

import (
	"context"
	"sharebasket/core"
	"sharebasket/domain/model"
	"sharebasket/domain/repository"
	"sharebasket/usecase/input"
	"sharebasket/usecase/output"
)

type (
	PersonalItem interface {
		Get(ctx context.Context, in input.GetPersonalItem) ([]output.GetPersonalItem, error)
		Create(ctx context.Context, in input.CreatePersonalItem) error
	}

	personalItemInteractor struct {
		accountRepo  repository.Account
		personalRepo repository.PersonalItem
		logger       core.Logger
	}
)

// [個人]買い物メモ一覧取得
func (p *personalItemInteractor) Get(ctx context.Context, in input.GetPersonalItem) ([]output.GetPersonalItem, error) {
	account, err := p.accountRepo.Get(in.UserID)
	if err != nil {
		p.logger.WithError(err).
			With("user_id", in.UserID).
			Error("failed to get account")
		return []output.GetPersonalItem{}, err
	}

	status, err := in.ParseShoppingStatus()
	if err != nil {
		p.logger.WithError(err).
			With("status", in.Status).
			Warn("invalid shopping status")
		return []output.GetPersonalItem{}, err
	}

	items, err := p.personalRepo.GetAll(account.ID, status)
	if err != nil {
		p.logger.WithError(err).
			Error("failed to get personal shopping item")
		return []output.GetPersonalItem{}, err
	}

	return output.NewGetPersonalItemOutputs(items), nil
}

// [個人]買い物メモ新規作成
func (p *personalItemInteractor) Create(ctx context.Context, in input.CreatePersonalItem) error {
	account, err := p.accountRepo.Get(in.UserID)
	if err != nil {
		p.logger.WithError(err).
			With("user_id", in.UserID).
			Error("failed to get account")
		return err
	}

	var status model.ShoppingStatus
	if in.Status != "" {
		status, err = model.ParseShoppingStatus(in.Status)
		if err != nil {
			p.logger.WithError(err).
				With("status", in.Status).
				Warn("invalid shopping status")
			return core.NewInvalidError(err)
		}
	}

	item, err := model.NewPersonalItem(in.Name, core.Ptr(status), in.CategoryID, account.ID)
	if err != nil {
		p.logger.WithError(err).
			With("name", in.Name).
			With("category_id", in.CategoryID).
			Warn("invalid personal item parameters")
		return core.NewInvalidError(err)
	}

	if err := p.personalRepo.Store(&item); err != nil {
		p.logger.WithError(err).
			With("item", item).
			Error("failed to store personal item")
		return err
	}

	return nil
}

func NewPersonalItem(a repository.Account, p repository.PersonalItem, l core.Logger) PersonalItem {
	return &personalItemInteractor{
		accountRepo:  a,
		personalRepo: p,
		logger:       l,
	}
}
