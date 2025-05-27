package auth

import "errors"

var (
	ErrLoginFailed        = errors.New("failed to login process")
	ErrInvalidInput       = errors.New("invalid input")
	ErrEmailExists        = errors.New("this email is already exists")
	ErrExpiredCode        = errors.New("expired code")
	ErrInvalidAccessToken = errors.New("invalid access token")
)
