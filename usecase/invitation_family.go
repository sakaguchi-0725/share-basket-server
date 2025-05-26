package usecase

import (
	"context"
	"errors"
	"sharebasket/core"
	"sharebasket/domain/repository"
	"sharebasket/domain/service"

	"github.com/google/uuid"
)

type (
	InvitationFamily interface {
		Execute(ctx context.Context, userID string) (string, error)
	}

	invitationFamily struct {
		accountRepo   repository.Account
		familyRepo    repository.Family
		familyService service.Family
	}
)

func (i *invitationFamily) Execute(ctx context.Context, userID string) (string, error) {
	account, err := i.accountRepo.Get(ctx, userID)
	if err != nil {
		return "", err
	}

	// 自身がオーナーの家族が存在するかチェック
	hasFamily, err := i.familyService.HasOwnedFamily(account.ID)
	if err != nil {
		return "", err
	}

	if !hasFamily {
		return "", core.NewInvalidError(errors.New("account does not have owner privileges for any family"))
	}

	family, err := i.familyRepo.GetOwnedFamily(account.ID)
	if err != nil {
		return "", err
	}

	// 招待可能上限数を超える場合
	if !family.CanInvite() {
		return "", core.NewInvalidError(errors.New("invitation limit reached for this family"))
	}

	token := uuid.NewString()
	err = i.familyRepo.Invitation(ctx, token, family.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func NewInvitationFamily(
	a repository.Account,
	f repository.Family,
	s service.Family,
) InvitationFamily {
	return &invitationFamily{
		accountRepo:   a,
		familyRepo:    f,
		familyService: s,
	}
}
