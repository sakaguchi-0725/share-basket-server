package core

type (
	ErrorCode int

	AppError struct {
		code ErrorCode
		error
	}
)

const (
	ErrBadRequest ErrorCode = iota + 1
	ErrUnauthorized
	ErrForbidden
	ErrEmailAlreadyExists
	ErrExpiredCode
	ErrExpiredToken
	ErrNotFound
)

func (e ErrorCode) String() string {
	switch e {
	case ErrBadRequest:
		return "BAD_REQUEST"
	case ErrEmailAlreadyExists:
		return "EMAIL_ALREADY_EXISTS"
	case ErrForbidden:
		return "FORBIDDEN"
	case ErrUnauthorized:
		return "UNAUTHORIZED"
	case ErrExpiredCode:
		return "EXPIRED_CODE"
	case ErrExpiredToken:
		return "EXPIRED_TOKEN"
	case ErrNotFound:
		return "DATA_NOT_FOUND"
	default:
		return "INTERNAL_SERVER_ERROR"
	}
}

func NewAppError(code ErrorCode, err error) *AppError {
	return &AppError{
		code:  code,
		error: err,
	}
}

func NewInvalidError(err error) *AppError {
	return &AppError{
		code:  ErrBadRequest,
		error: err,
	}
}

func (e *AppError) Error() string {
	return e.error.Error()
}

func (e *AppError) Code() ErrorCode {
	return e.code
}
