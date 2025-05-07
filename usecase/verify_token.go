//go:generate mockgen -destination=../test/mock/usecase/verify_token_input.go . VerifyTokenInputPort
package usecase

import (
	"context"
	"errors"
	"share-basket-server/core/apperr"
	"share-basket-server/core/logger"
	"share-basket-server/domain"
)

type (
	VerifyTokenInputPort interface {
		Execute(ctx context.Context, token string) (string, error)
	}

	verifyTokenInteractor struct {
		authenticator domain.Authenticator
		userRepo      domain.UserRepository
		logger        logger.Logger
	}
)

// tokenを検証し、認証されたUserIDを返す
func (v *verifyTokenInteractor) Execute(ctx context.Context, token string) (string, error) {
	email, err := v.authenticator.VerifyToken(ctx, token)
	if err != nil {
		if errors.Is(err, apperr.ErrInvalidToken) || errors.Is(err, apperr.ErrTokenExpired) {
			v.logger.
				With("token", token).
				Info("invalid token")
			return "", apperr.New(apperr.ErrUnauthorized, err)
		}
		v.logger.
			With("token", token).
			Error("failed to verify token")
		return "", err
	}

	user, err := v.userRepo.GetByEmail(email)
	if err != nil {
		if errors.Is(err, apperr.ErrDataNotFound) {
			v.logger.
				With("email", email).
				Info("user not found")
			return "", apperr.New(apperr.ErrUnauthorized, err)
		}

		v.logger.
			With("email", email).
			Error("failed to get user")
		return "", err
	}

	return user.ID.String(), nil
}

func NewVerifyTokenInteractor(
	authenticator domain.Authenticator,
	userRepo domain.UserRepository,
	logger logger.Logger,
) VerifyTokenInputPort {
	return &verifyTokenInteractor{authenticator, userRepo, logger}
}
