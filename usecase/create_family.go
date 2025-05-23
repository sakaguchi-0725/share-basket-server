package usecase

import (
	"context"
	"errors"
	"sharebasket/core"
	"sharebasket/domain/model"
	"sharebasket/domain/repository"
	"sharebasket/domain/service"
)

type (
	CreateFamily interface {
		Execute(ctx context.Context, in CreateFamilyInput) error
	}

	CreateFamilyInput struct {
		Name   string
		UserID string
	}

	createFamily struct {
		accountRepo   repository.Account
		familyRepo    repository.Family
		familyService service.Family
	}
)

// 新しい家族を作成します。
// 指定されたユーザーIDに対応するアカウントを検証し、新しい家族を作成してリポジトリに保存します。
// アカウントが見つからない場合は認証エラーを、入力が無効な場合は不正エラーを返します。
// 成功した場合はnilを返します。
func (c *createFamily) Execute(ctx context.Context, in CreateFamilyInput) error {
	account, err := c.accountRepo.Get(in.UserID)
	if err != nil {
		return err
	}

	// すでに家族オーナー、または家族メンバーではないか判定
	hasFamily, err := c.familyService.HasFamily(account.ID)
	if err != nil {
		return err
	}

	if hasFamily {
		return core.NewInvalidError(errors.New("account already has family"))
	}

	id := model.NewFamilyID()

	family, err := model.NewFamily(id, in.Name, account)
	if err != nil {
		return core.NewInvalidError(err)
	}

	err = c.familyRepo.Store(&family)
	if err != nil {
		return err
	}

	return nil
}

func NewCreateFamily(a repository.Account, f repository.Family, s service.Family) CreateFamily {
	return &createFamily{
		accountRepo:   a,
		familyRepo:    f,
		familyService: s,
	}
}
