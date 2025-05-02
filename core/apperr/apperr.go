package apperr

type AppError struct {
	code ErrorCode
	error
}

func New(code ErrorCode, err error) *AppError {
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

func (err *AppError) Code() ErrorCode {
	return err.code
}
