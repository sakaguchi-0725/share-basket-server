package usecase

import (
	"context"
	"errors"
	"share-basket-server/core/apperr"
	"share-basket-server/domain"
)

type (
	// Tokenを検証し、ログイン済みのUserIDを返す
	VerifyTokenInputPort interface {
		Execute(ctx context.Context, token string) (string, error)
	}

	verifyTokenInteractor struct {
		authenticator domain.Authenticator
		userRepo      domain.UserRepository
	}
)

func (v *verifyTokenInteractor) Execute(ctx context.Context, token string) (string, error) {
	email, err := v.authenticator.VerifyToken(ctx, token)
	if err != nil {
		if errors.Is(err, apperr.ErrInvalidToken) || errors.Is(err, apperr.ErrTokenExpired) {
			return "", apperr.New(apperr.ErrUnauthorized, err)
		}

		return "", err
	}

	user, err := v.userRepo.GetByEmail(email)
	if err != nil {
		if errors.Is(err, apperr.ErrDataNotFound) {
			return "", apperr.New(apperr.ErrUnauthorized, err)
		}

		return "", err
	}

	return user.ID.String(), nil
}

func NewVerifyTokenInteractor(
	authenticator domain.Authenticator,
	userRepo domain.UserRepository,
) VerifyTokenInputPort {
	return &verifyTokenInteractor{authenticator, userRepo}
}
