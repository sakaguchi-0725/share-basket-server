//go:generate mockgen -destination=../test/mock/usecase/signup_confirm_input.go . SignUpConfirmInputPort
//go:generate mockgen -destination=../test/mock/usecase/signup_confirm_output.go . SignUpConfirmOutputPort
package usecase

import (
	"context"
	"errors"
	"share-basket-server/core/apperr"
	"share-basket-server/core/logger"
	"share-basket-server/domain"
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
		logger        logger.Logger
	}
)

func (s *signUpConfirmInteractor) Execute(
	ctx context.Context, input SignUpConfirmInput, output SignUpConfirmOutputPort,
) error {
	err := s.authenticator.SignUpConfirm(ctx, input.Email, input.ConfirmationCode)

	if err != nil {
		if errors.Is(err, apperr.ErrInvalidData) {
			s.logger.
				With("email", input.Email).
				With("error", err).
				Info("invalid input")
			return apperr.NewInvalidError(err)
		}

		if errors.Is(err, apperr.ErrExpiredCodeException) {
			s.logger.
				With("code", input.ConfirmationCode).
				With("error", err).
				Info("expired code")
			return apperr.New(apperr.ErrExpiredCode, err)
		}

		s.logger.
			With("error", err).
			Error("failed to sigin up confirm")
		return err
	}

	return output.Render(ctx)
}

func NewSignUpConfirmInteractor(
	authenticator domain.Authenticator,
	logger logger.Logger,
) SignUpConfirmInputPort {
	return &signUpConfirmInteractor{authenticator, logger}
}
