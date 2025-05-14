package usecase

import (
	"context"
	"errors"
	"sharebasket/core"
	"sharebasket/domain/repository"
)

var (
	ErrInvalidAccessToken = errors.New("invalid access token")
	ErrExpiredAccessToken = errors.New("access token is expired")
)

type (
	VerifyToken interface {
		Execute(ctx context.Context, token string) (string, error)
	}

	verifyToken struct {
		auth   repository.Authenticator
		logger core.Logger
	}
)

func (v *verifyToken) Execute(ctx context.Context, token string) (string, error) {
	userID, err := v.auth.VerifyToken(ctx, token)
	if err != nil {
		if errors.Is(err, ErrInvalidAccessToken) {
			v.logger.WithError(err).
				With("access_token", token).
				Warn("invalid access token")
			return "", core.NewAppError(core.ErrUnauthorized, err)
		}
		if errors.Is(err, ErrExpiredAccessToken) {
			v.logger.WithError(err).
				With("access_token", token).
				Warn("access token is expired")
			return "", core.NewAppError(core.ErrExpiredToken, err)
		}

		v.logger.WithError(err).
			With("access_token", token).
			Error("failed to verify access token")
		return "", err
	}

	return userID, nil
}

func NewVerifyToken(a repository.Authenticator, l core.Logger) VerifyToken {
	return &verifyToken{
		auth:   a,
		logger: l,
	}
}
