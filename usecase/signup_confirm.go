package usecase

import (
	"context"
	"sharebasket/domain/repository"
)

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
		return err
	}

	return nil
}

func NewSignUpConfirm(a repository.Authenticator) SignUpConfirm {
	return &signUpConfirm{
		auth: a,
	}
}
