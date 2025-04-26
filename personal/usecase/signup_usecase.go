package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"share-basket-server/core/apperr"
	"share-basket-server/personal/domain/model"
	"share-basket-server/personal/domain/repository"
	"share-basket-server/personal/domain/service"
	"share-basket-server/personal/usecase/input"
	"share-basket-server/personal/usecase/output"
)

type signup struct {
	authenticator repository.Authenticator
	userRepo      repository.User
	accountRepo   repository.Account
	userService   service.UserService
	transaction   repository.Transaction
}

func (s *signup) Execute(ctx context.Context, in input.SignUpInput, out output.SignUpOutputPort) error {
	emailAvailable, err := s.userService.IsEmailAvailable(in.Email)
	if err != nil {
		slog.Error("isEmailAvailable failed", slog.String("error", err.Error()))
		return err
	}

	if !emailAvailable {
		slog.Info("email is not available")
		return apperr.NewInvalidError(fmt.Errorf("email is not available: %s", in.Email))
	}

	cognitoUID, err := s.authenticator.SignUp(ctx, in.Email, in.Password)
	if err != nil {
		if errors.Is(err, apperr.ErrDuplicatedKey) || errors.Is(err, apperr.ErrInvalidData) {
			slog.Info("sign up failed", slog.String("error", err.Error()))
			return apperr.NewInvalidError(err)
		}
		slog.Error("sign up failed", slog.String("error", err.Error()))
		return err
	}

	err = s.transaction.Run(ctx, func(ctx context.Context) error {
		user, err := model.NewUser(model.GenUserID(), cognitoUID, in.Email)
		if err != nil {
			return apperr.NewInvalidError(err)
		}

		if err := s.userRepo.Store(&user); err != nil {
			return err
		}

		acc, err := model.NewAccount(model.GenAccountID(), user.ID, in.Name)
		if err != nil {
			return apperr.NewInvalidError(err)
		}

		if err := s.accountRepo.Store(&acc); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		if err := s.authenticator.DeleteUser(ctx, in.Email); err != nil {
			return err
		}
		return err
	}

	return out.Render(ctx)
}

func NewSignUpUseCase(
	authenticator repository.Authenticator,
	userRepo repository.User,
	accountRepo repository.Account,
	userService service.UserService,
	transaction repository.Transaction,
) input.SignUpInputPort {
	return &signup{
		authenticator,
		userRepo,
		accountRepo,
		userService,
		transaction,
	}
}
