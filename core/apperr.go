package core

import (
	"fmt"
	"path/filepath"
	"runtime"
)

type AppError struct {
	code    ErrorCode
	file    string
	line    int
	message string
	error
}

func NewAppError(code ErrorCode, err error) *AppError {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "unknown"
		line = 0
	}

	return &AppError{
		code:    code,
		error:   err,
		file:    filepath.Base(file),
		line:    line,
		message: code.DefaultMessage(),
	}
}

func NewInvalidError(err error) *AppError {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "unknown"
		line = 0
	}

	return &AppError{
		code:    ErrBadRequest,
		error:   err,
		file:    filepath.Base(file),
		line:    line,
		message: ErrBadRequest.DefaultMessage(),
	}
}

func NewUnauthorizedError(err error) *AppError {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "unknown"
		line = 0
	}

	return &AppError{
		code:    ErrUnauthorized,
		error:   err,
		file:    filepath.Base(file),
		line:    line,
		message: ErrBadRequest.DefaultMessage(),
	}
}

// カスタムメッセージを設定
func (e *AppError) WithMessage(message string) *AppError {
	e.message = message
	return e
}

func (e *AppError) Unwrap() error {
	return e.error
}

func (e *AppError) Is(target error) bool {
	if t, ok := target.(*AppError); ok {
		return e.code == t.code
	}
	return false
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s at %s:%d: %s", e.code.String(), e.file, e.line, e.error.Error())
}

func (e *AppError) Code() ErrorCode {
	return e.code
}

func (e *AppError) StackTrace() string {
	return fmt.Sprintf("%s:%d", e.file, e.line)
}

func (e *AppError) Message() string {
	return e.message
}
