package usecase

import (
	"context"
	"sharebasket/domain/repository"
)

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
