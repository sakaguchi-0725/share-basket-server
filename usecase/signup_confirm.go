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
		auth repository.Authenticator
	}
)

func (s *signUpConfirm) Execute(ctx context.Context, in SignUpConfirmInput) error {
	err := s.auth.SignUpConfirm(ctx, in.Email, in.ConfirmationCode)
	if err != nil {
		if errors.Is(err, core.ErrInvalidData) {
			return core.NewInvalidError(err)
		}

		if errors.Is(err, ErrExpiredConfirmationCode) {
			return core.NewAppError(core.ErrExpiredCode, err)
		}

		return err
	}

	return nil
}

func NewSignUpConfirm(a repository.Authenticator) SignUpConfirm {
	return &signUpConfirm{auth: a}
}
