package core

import "errors"

var (
	ErrInvalidData  = errors.New("invalid data")
	ErrDataNotFound = errors.New("data not found")
)
