package repository

import "context"

type Authenticator interface {
	SignUp(ctx context.Context, email, password string) (string, error)
	SignUpConfirm(ctx context.Context, email, confirmationCode string) error
	Login(ctx context.Context, email, password string) (accessToken, refreshToken string, err error)
	VerifyToken(ctx context.Context, token string) (string, error)
	DeleteUser(ctx context.Context, email string) error
}
