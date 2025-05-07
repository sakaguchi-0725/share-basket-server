//go:generate mockgen -destination=../test/mock/usecase/signup_input.go . SignUpInputPort
//go:generate mockgen -destination=../test/mock/usecase/signup_output.go . SignUpOutputPort
package usecase

import (
	"context"
	"errors"
	"fmt"
	"share-basket-server/core/apperr"
	"share-basket-server/core/logger"
	"share-basket-server/domain"
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
		logger        logger.Logger
	}
)

func (s *signupInteractor) Execute(ctx context.Context, input SignUpInput, output SignUpOutputPort) error {
	emailAvailable, err := s.userService.IsEmailAvailable(input.Email)
	if err != nil {
		s.logger.
			With("error", err).
			Error("failed to check email available")
		return err
	}

	if !emailAvailable {
		s.logger.
			With("error", err).
			With("email", input.Email).
			Info("email is not available")
		return apperr.NewInvalidError(fmt.Errorf("email is not available: %s", input.Email))
	}

	cognitoUID, err := s.authenticator.SignUp(ctx, input.Email, input.Password)
	if err != nil {
		if errors.Is(err, apperr.ErrDuplicatedKey) || errors.Is(err, apperr.ErrInvalidData) {
			s.logger.
				With("error", err).
				Info("invalid input")
			return apperr.NewInvalidError(err)
		}

		s.logger.
			With("error", err).
			Error("failed to sign up")
		return err
	}

	err = s.transaction.Run(ctx, func(ctx context.Context) error {
		user, err := domain.NewUser(domain.NewUserID(), cognitoUID, input.Email)
		if err != nil {
			s.logger.
				With("error", err).
				Info("failed to NewUser")
			return apperr.NewInvalidError(err)
		}

		if err := s.userRepo.Store(&user); err != nil {
			s.logger.
				With("error", err).
				Error("failed to store user")
			return err
		}

		acc, err := domain.NewAccount(domain.NewAccountID(), user.ID, input.Name)
		if err != nil {
			s.logger.
				With("error", err).
				Error("failed to NewAccount")
			return apperr.NewInvalidError(err)
		}

		if err := s.accountRepo.Store(&acc); err != nil {
			s.logger.
				With("error", err).
				Error("failed to store account")
			return err
		}

		return nil
	})

	if err != nil {
		if err := s.authenticator.DeleteUser(ctx, input.Email); err != nil {
			s.logger.
				With("error", err).
				Error("failed to delete user")
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
	logger logger.Logger,
) SignUpInputPort {
	return &signupInteractor{
		authenticator,
		userRepo,
		accountRepo,
		userService,
		transaction,
		logger,
	}
}
