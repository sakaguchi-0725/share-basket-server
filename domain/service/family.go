package service

import (
	"sharebasket/domain/model"
	"sharebasket/domain/repository"
)

type (
	Family interface {
		HasFamily(accountID model.AccountID) (bool, error)
	}

	familyService struct {
		repo repository.Family
	}
)

// HasFamily は指定されたアカウントIDが家族のオーナーまたはメンバーとして存在するかを確認する
func (f *familyService) HasFamily(accountID model.AccountID) (bool, error) {
	return f.repo.HasFamily(accountID)
}

func NewFamily(r repository.Family) Family {
	return &familyService{repo: r}
}
