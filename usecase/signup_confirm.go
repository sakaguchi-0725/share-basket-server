package usecase

import (
	"context"
	"errors"
	"sharebasket/core"
	"sharebasket/domain/repository"
)

var ErrExpiredConfirmationCode = errors.New("confirmation code is expired")

type (
	SignUpConfirm interface {
		Execute(ctx context.Context, in SignUpConfirmInput) error
	}

	SignUpConfirmInput struct {
		Email            string
		ConfirmationCode string
	}

	signUpConfirm struct {
		auth   repository.Authenticator
		logger core.Logger
	}
)

func (s *signUpConfirm) Execute(ctx context.Context, in SignUpConfirmInput) error {
	err := s.auth.SignUpConfirm(ctx, in.Email, in.ConfirmationCode)
	if err != nil {
		if errors.Is(err, core.ErrInvalidData) {
			s.logger.WithError(err).
				With("email", in.Email).
				With("confirmationCode", in.ConfirmationCode).
				Warn("invalid sign up confirm input")
			return core.NewInvalidError(err)
		}

		if errors.Is(err, ErrExpiredConfirmationCode) {
			s.logger.WithError(err).
				With("email", in.Email).
				With("confirmationCode", in.ConfirmationCode).
				Warn("confirmation code is expired")
			return core.NewAppError(core.ErrExpiredCode, err)
		}

		s.logger.WithError(err).
			With("email", in.Email).
			With("confirmationCode", in.ConfirmationCode).
			Error("failed to signup confirmation")
		return err
	}

	return nil
}

func NewSignUpConfirm(a repository.Authenticator, l core.Logger) SignUpConfirm {
	return &signUpConfirm{
		auth:   a,
		logger: l,
	}
}
