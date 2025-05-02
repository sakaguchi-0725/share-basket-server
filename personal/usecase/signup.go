package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"share-basket-server/core/apperr"
	"share-basket-server/personal/domain"
)

type (
	SignUpInputPort interface {
		Execute(ctx context.Context, input SignUpInput, output SignUpOutputPort) error
	}

	SignUpInput struct {
		Name     string
		Email    string
		Password string
	}

	SignUpOutputPort interface {
		Render(ctx context.Context) error
	}
	signupInteractor struct {
		authenticator domain.Authenticator
		userRepo      domain.UserRepository
		accountRepo   domain.AccountRepository
		userService   domain.UserService
		transaction   domain.Transaction
	}
)

func (s *signupInteractor) Execute(ctx context.Context, input SignUpInput, output SignUpOutputPort) error {
	emailAvailable, err := s.userService.IsEmailAvailable(input.Email)
	if err != nil {
		slog.Error("isEmailAvailable failed", slog.String("error", err.Error()))
		return err
	}

	if !emailAvailable {
		slog.Info("email is not available")
		return apperr.NewInvalidError(fmt.Errorf("email is not available: %s", input.Email))
	}

	cognitoUID, err := s.authenticator.SignUp(ctx, input.Email, input.Password)
	if err != nil {
		if errors.Is(err, apperr.ErrDuplicatedKey) || errors.Is(err, apperr.ErrInvalidData) {
			slog.Info("sign up failed", slog.String("error", err.Error()))
			return apperr.NewInvalidError(err)
		}
		slog.Error("sign up failed", slog.String("error", err.Error()))
		return err
	}

	err = s.transaction.Run(ctx, func(ctx context.Context) error {
		user, err := domain.NewUser(domain.NewUserID(), cognitoUID, input.Email)
		if err != nil {
			return apperr.NewInvalidError(err)
		}

		if err := s.userRepo.Store(&user); err != nil {
			return err
		}

		acc, err := domain.NewAccount(domain.NewAccountID(), user.ID, input.Name)
		if err != nil {
			return apperr.NewInvalidError(err)
		}

		if err := s.accountRepo.Store(&acc); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		if err := s.authenticator.DeleteUser(ctx, input.Email); err != nil {
			return err
		}
		return err
	}

	return output.Render(ctx)
}

func NewSignUpInteractor(
	authenticator domain.Authenticator,
	userRepo domain.UserRepository,
	accountRepo domain.AccountRepository,
	userService domain.UserService,
	transaction domain.Transaction,
) SignUpInputPort {
	return &signupInteractor{
		authenticator,
		userRepo,
		accountRepo,
		userService,
		transaction,
	}
}
