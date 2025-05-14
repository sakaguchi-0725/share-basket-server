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

var ErrEmailAlreadyExists = errors.New("email is already exists")

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
		logger      core.Logger
	}
)

// Execute は、新規ユーザーのサインアップ処理を実行します。
// メールアドレスの重複チェック、認証プロバイダでのユーザー作成、
// ユーザーとアカウントの永続化を行います。
// サインアップ中にエラーが発生した場合、トランザクションをロールバックし、
// 作成途中のユーザーを削除します。
func (s *signUp) Execute(ctx context.Context, in SignUpInput) (err error) {
	available, err := s.userService.IsEmailAvailable(in.Email)
	if err != nil {
		s.logger.WithError(err).
			With("email", in.Email).
			Error("failed to check email availability")
		return fmt.Errorf("failed to check email availability: %w", err)
	}

	if !available {
		s.logger.With("email", in.Email).Warn("email is exists")
		return core.NewAppError(core.ErrEmailAlreadyExists, ErrEmailAlreadyExists)
	}

	id, err := s.auth.SignUp(ctx, in.Email, in.Password)
	if err != nil {
		if errors.Is(err, ErrEmailAlreadyExists) {
			s.logger.WithError(err).
				With("email", in.Email).
				Warn("email is exists")
			return core.NewAppError(core.ErrEmailAlreadyExists, err)
		}

		if errors.Is(err, core.ErrInvalidData) {
			s.logger.WithError(err).
				With("name", in.Name).
				With("email", in.Email).
				Warn("invalid sign up input")
			return core.NewInvalidError(err)
		}

		s.logger.WithError(err).
			With("name", in.Name).
			With("email", in.Email).
			Error("failed to sign up")
		return fmt.Errorf("failed to sign up: %w", err)
	}

	err = s.tx.WithTx(ctx, func(ctx context.Context) error {
		user, err := model.NewUser(id, in.Email)
		if err != nil {
			s.logger.WithError(err).
				With("id", id).
				With("email", in.Email).
				Warn("failed to new user model")
			return core.NewInvalidError(err)
		}

		if err = s.userRepo.Store(ctx, &user); err != nil {
			s.logger.WithError(err).Error("failed to store user")
			return fmt.Errorf("failed to store user: %w", err)
		}

		accID := model.NewAccountID()
		acc, err := model.NewAccount(accID, in.Name, user.ID)
		if err != nil {
			s.logger.WithError(err).
				With("id", accID.String()).
				With("name", in.Name).
				With("user_id", user.ID).
				Warn("failed to new account model")
			return core.NewInvalidError(err)
		}

		if err = s.accountRepo.Store(ctx, &acc); err != nil {
			s.logger.WithError(err).Error("failed to store account")
			return fmt.Errorf("failed to store account: %w", err)
		}

		return nil
	})

	if err != nil {
		if delErr := s.auth.DeleteUser(ctx, in.Email); delErr != nil {
			s.logger.WithError(err).
				With("email", in.Email).
				Error("failed to delete cognito user")
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
	l core.Logger,
) SignUp {
	return &signUp{
		auth:        a,
		userRepo:    ur,
		accountRepo: ar,
		userService: us,
		tx:          tx,
		logger:      l,
	}
}
