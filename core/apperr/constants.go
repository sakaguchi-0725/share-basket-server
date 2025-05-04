package apperr

import "errors"

var (
	ErrInvalidData          = errors.New("invalid data")
	ErrDataNotFound         = errors.New("data not found")
	ErrDuplicatedKey        = errors.New("duplicate key already exists")
	ErrExpiredCodeException = errors.New("code has expired")
	ErrUnauthenticated      = errors.New("unauthenticated")
	ErrInvalidToken         = errors.New("invalid token")
	ErrTokenExpired         = errors.New("token is expired")
)
