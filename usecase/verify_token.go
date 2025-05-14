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
		auth repository.Authenticator
	}
)

func (v *verifyToken) Execute(ctx context.Context, token string) (string, error) {
	userID, err := v.auth.VerifyToken(ctx, token)
	if err != nil {
		if errors.Is(err, ErrInvalidAccessToken) {
			return "", core.NewAppError(core.ErrUnauthorized, err)
		}
		if errors.Is(err, ErrExpiredAccessToken) {
			return "", core.NewAppError(core.ErrExpiredToken, err)
		}

		return "", err
	}

	return userID, nil
}

func NewVerifyToken(a repository.Authenticator) VerifyToken {
	return &verifyToken{auth: a}
}
