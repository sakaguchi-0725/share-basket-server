package usecase

import (
	"context"
	"errors"
	"share-basket-server/core/apperr"
	"share-basket-server/personal/domain/repository"
	"share-basket-server/personal/usecase/input"
	"share-basket-server/personal/usecase/output"
)

type loginUseCase struct {
	authenticator repository.Authenticator
}

func (l *loginUseCase) Execute(ctx context.Context, in input.LoginInput, out output.LoginOutputPort) error {
	accessToken, err := l.authenticator.Login(ctx, in.Email, in.Password)

	if err != nil {
		if errors.Is(err, apperr.ErrUnauthenticated) {
			return apperr.New(apperr.ErrUnauthorized, err)
		}
		if errors.Is(err, apperr.ErrDataNotFound) {
			return apperr.NewInvalidError(err)
		}

		return err
	}

	return out.Render(ctx, accessToken)
}

func NewLoginUseCase(authenticator repository.Authenticator) input.LoginInputPort {
	return &loginUseCase{authenticator}
}
