package usecase

import (
	"context"
	"errors"
	"fmt"
	"sharebasket/core"
	"sharebasket/domain/model"
	"sharebasket/domain/repository"
	"sharebasket/domain/service"
)

type (
	SignUp interface {
		Execute(ctx context.Context, in SignUpInput) error
	}

	SignUpInput struct {
		Name     string
		Email    string
		Password string
	}

	signUp struct {
		auth        repository.Authenticator
		userRepo    repository.User
		accountRepo repository.Account
		userService service.User
		tx          repository.Transaction
	}
)

func (s *signUp) Execute(ctx context.Context, in SignUpInput) (err error) {
	// メールアドレスが使用可能がチェック
	available, err := s.userService.IsEmailAvailable(ctx, in.Email)
	if err != nil {
		return fmt.Errorf("failed to check email availability: %w", err)
	}

	if !available {
		return core.NewAppError(
			core.ErrEmailAlreadyExists, errors.New("this email is already exists"),
		).WithMessage("このメールアドレスは使用できません")
	}

	id, err := s.auth.SignUp(ctx, in.Email, in.Password)
	if err != nil {
		return err
	}

	// トランザクション
	err = s.tx.WithTx(ctx, func(ctx context.Context) error {
		user, err := model.NewUser(id, in.Email)
		if err != nil {
			return err
		}

		if err = s.userRepo.Store(ctx, &user); err != nil {
			return err
		}

		accID := model.NewAccountID()
		acc, err := model.NewAccount(accID, in.Name, user.ID)
		if err != nil {
			return err
		}

		if err = s.accountRepo.Store(ctx, &acc); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		// 登録に失敗した場合、Cognitoのユーザーを削除する
		if delErr := s.auth.DeleteUser(ctx, in.Email); delErr != nil {
			return fmt.Errorf("failed to delete cognito user after transaction error: %w", delErr)
		}
		return err
	}

	return nil
}

func NewSignUp(
	a repository.Authenticator,
	ur repository.User,
	ar repository.Account,
	us service.User,
	tx repository.Transaction,
) SignUp {
	return &signUp{
		auth:        a,
		userRepo:    ur,
		accountRepo: ar,
		userService: us,
		tx:          tx,
	}
}
