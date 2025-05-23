package usecase

import (
	"context"
	"errors"
	"sharebasket/core"
	"sharebasket/domain/repository"
)

var ErrLoginFailed = errors.New("failed to login process")

type (
	Login interface {
		Execute(ctx context.Context, in LoginInput) (LoginOutput, error)
	}

	LoginInput struct {
		Email    string
		Password string
	}

	LoginOutput struct {
		AccessToken  string
		RefreshToken string
	}

	login struct {
		auth repository.Authenticator
	}
)

func (l *login) Execute(ctx context.Context, in LoginInput) (LoginOutput, error) {
	accessToken, refreshToken, err := l.auth.Login(ctx, in.Email, in.Password)
	if err != nil {
		if errors.Is(err, ErrLoginFailed) || errors.Is(err, core.ErrInvalidData) {
			return LoginOutput{}, core.NewAppError(core.ErrUnauthorized, err)
		}

		return LoginOutput{}, err
	}

	return LoginOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func NewLogin(a repository.Authenticator) Login {
	return &login{
		auth: a,
	}
}
