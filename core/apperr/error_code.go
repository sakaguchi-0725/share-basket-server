package apperr

type ErrorCode int

const (
	ErrUnauthorized ErrorCode = iota + 1
	ErrBadRequest
	ErrNotFound
)

func (code ErrorCode) String() string {
	switch code {
	case ErrUnauthorized:
		return "Unauthorized"
	case ErrBadRequest:
		return "BadRequest"
	case ErrNotFound:
		return "NotFound"
	default:
		return "InternalServer"
	}
}
