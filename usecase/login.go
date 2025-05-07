//go:generate mockgen -destination=../test/mock/usecase/login_input.go . LoginInputPort
//go:generate mockgen -destination=../test/mock/usecase/login_output.go . LoginOutputPort
package usecase

import (
	"context"
	"errors"
	"share-basket-server/core/apperr"
	"share-basket-server/core/logger"
	"share-basket-server/domain"
)

type (
	LoginInputPort interface {
		Execute(ctx context.Context, input LoginInput, output LoginOutputPort) error
	}

	LoginInput struct {
		Email    string
		Password string
	}

	LoginOutputPort interface {
		Render(ctx context.Context, token string) error
	}

	loginInteractor struct {
		authenticator domain.Authenticator
		logger        logger.Logger
	}
)

func (l *loginInteractor) Execute(ctx context.Context, input LoginInput, output LoginOutputPort) error {
	accessToken, err := l.authenticator.Login(ctx, input.Email, input.Password)

	if err != nil {
		if errors.Is(err, apperr.ErrUnauthenticated) {
			l.logger.
				With("email", input.Email).
				With("error", err).
				Info("unautorized")
			return apperr.New(apperr.ErrUnauthorized, err)
		}
		if errors.Is(err, apperr.ErrDataNotFound) {
			l.logger.
				With("email", input.Email).
				With("error", err).
				Info("user not found")
			return apperr.NewInvalidError(err)
		}

		l.logger.
			With("email", input.Email).
			With("error", err).
			Error("failed to login")
		return err
	}

	return output.Render(ctx, accessToken)
}

func NewLoginInteractor(authenticator domain.Authenticator, logger logger.Logger) LoginInputPort {
	return &loginInteractor{authenticator, logger}
}
