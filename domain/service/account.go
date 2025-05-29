package service

import (
	"context"
	"sharebasket/domain/model"
	"sharebasket/domain/repository"
)

type (
	Account interface {
		HasFamily(ctx context.Context, id model.AccountID) (bool, error)
		HasOwnedFamily(ctx context.Context, id model.AccountID) (bool, error)
		HasMembership(ctx context.Context, accountID model.AccountID, familyID model.FamilyID) (bool, error)
	}

	accountService struct {
		repo repository.Family
	}
)

// 家族に参加しているか確認する
func (a *accountService) HasFamily(ctx context.Context, id model.AccountID) (bool, error) {
	return a.repo.HasFamily(ctx, id)
}

// 自身がオーナーの家族があるか確認する
func (a *accountService) HasOwnedFamily(ctx context.Context, id model.AccountID) (bool, error) {
	return a.repo.HasOwnedFamily(ctx, id)
}

// アカウントが特定の家族のメンバーシップを持っているか確認する
func (a *accountService) HasMembership(ctx context.Context, accountID model.AccountID, familyID model.FamilyID) (bool, error) {
	return a.repo.HasMembership(ctx, accountID, familyID)
}

func NewAccount(r repository.Family) Account {
	return &accountService{repo: r}
}
