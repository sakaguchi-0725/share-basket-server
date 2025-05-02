package usecase

import (
	"context"
	"errors"
	"share-basket-server/core/apperr"
	"share-basket-server/personal/domain"
)

type (
	SignUpConfirmInputPort interface {
		Execute(ctx context.Context, input SignUpConfirmInput, output SignUpConfirmOutputPort) error
	}

	SignUpConfirmInput struct {
		Email            string
		ConfirmationCode string
	}

	SignUpConfirmOutputPort interface {
		Render(ctx context.Context) error
	}

	signUpConfirmInteractor struct {
		authenticator domain.Authenticator
	}
)

func (s *signUpConfirmInteractor) Execute(
	ctx context.Context, input SignUpConfirmInput, output SignUpConfirmOutputPort,
) error {
	err := s.authenticator.SignUpConfirm(ctx, input.Email, input.ConfirmationCode)

	if err != nil {
		if errors.Is(err, apperr.ErrInvalidData) {
			return apperr.NewInvalidError(err)
		}
		if errors.Is(err, apperr.ErrExpiredCodeException) {
			return apperr.New(apperr.ErrExpiredCode, err)
		}
		return err
	}

	return output.Render(ctx)
}

func NewSignUpConfirmInteractor(authenticator domain.Authenticator) SignUpConfirmInputPort {
	return &signUpConfirmInteractor{authenticator}
}
