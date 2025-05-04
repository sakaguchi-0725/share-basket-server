//go:generate mockgen -destination=../mock/usecase/login_input.go . LoginInputPort
//go:generate mockgen -destination=../mock/usecase/login_output.go . LoginOutputPort
package usecase

import (
	"context"
	"errors"
	"share-basket-server/core/apperr"
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
	}
)

func (l *loginInteractor) Execute(ctx context.Context, input LoginInput, output LoginOutputPort) error {
	accessToken, err := l.authenticator.Login(ctx, input.Email, input.Password)

	if err != nil {
		if errors.Is(err, apperr.ErrUnauthenticated) {
			return apperr.New(apperr.ErrUnauthorized, err)
		}
		if errors.Is(err, apperr.ErrDataNotFound) {
			return apperr.NewInvalidError(err)
		}

		return err
	}

	return output.Render(ctx, accessToken)
}

func NewLoginInteractor(authenticator domain.Authenticator) LoginInputPort {
	return &loginInteractor{authenticator}
}
