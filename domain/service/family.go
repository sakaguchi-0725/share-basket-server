package service

import (
	"context"
	"sharebasket/domain/model"
	"sharebasket/domain/repository"
)

type (
	Family interface {
		HasFamily(ctx context.Context, accountID model.AccountID) (bool, error)
		HasOwnedFamily(ctx context.Context, accountID model.AccountID) (bool, error)
	}

	familyService struct {
		repo repository.Family
	}
)

// 自身がオーナーの家族があるか確認する
func (f *familyService) HasOwnedFamily(ctx context.Context, accountID model.AccountID) (bool, error) {
	return f.repo.HasOwnedFamily(ctx, accountID)
}

// HasFamily は指定されたアカウントIDが家族のオーナーまたはメンバーとして存在するかを確認する
func (f *familyService) HasFamily(ctx context.Context, accountID model.AccountID) (bool, error) {
	return f.repo.HasFamily(ctx, accountID)
}

func NewFamily(r repository.Family) Family {
	return &familyService{repo: r}
}
