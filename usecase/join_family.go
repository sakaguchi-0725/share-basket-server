package usecase

import (
	"context"
	"sharebasket/core"
	"sharebasket/domain/repository"
)

type (
	JoinFamily interface {
		Execute(ctx context.Context, in JoinFamilyInput) error
	}

	JoinFamilyInput struct {
		UserID string
		Token  string
	}

	joinFamily struct {
		accountRepo repository.Account
		familyRepo  repository.Family
	}
)

func (j *joinFamily) Execute(ctx context.Context, in JoinFamilyInput) error {
	account, err := j.accountRepo.Get(in.UserID)
	if err != nil {
		return err
	}

	family, err := j.familyRepo.GetByToken(ctx, in.Token)
	if err != nil {
		return err
	}

	if err := family.Join(account.ID); err != nil {
		return core.NewInvalidError(err)
	}

	if err := j.familyRepo.Store(&family); err != nil {
		return err
	}

	return nil
}

func NewJoinFamily(a repository.Account, f repository.Family) JoinFamily {
	return &joinFamily{
		accountRepo: a,
		familyRepo:  f,
	}
}
