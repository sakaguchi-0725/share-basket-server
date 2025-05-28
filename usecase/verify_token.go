package usecase

import (
	"context"
	"sharebasket/domain/repository"
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
		return "", err
	}

	return userID, nil
}

func NewVerifyToken(a repository.Authenticator) VerifyToken {
	return &verifyToken{
		auth: a,
	}
}
