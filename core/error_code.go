package core

type ErrorCode int

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

// エラーコードに対応するデフォルトメッセージを返す
func (e ErrorCode) DefaultMessage() string {
	switch e {
	case ErrBadRequest:
		return "リクエストが不正です"
	case ErrEmailAlreadyExists:
		return "このメールアドレスは使用できません"
	case ErrForbidden:
		return "アクセス権限がありません"
	case ErrUnauthorized:
		return "認証が必要です"
	case ErrExpiredCode:
		return "コードの有効期限が切れています"
	case ErrExpiredToken:
		return "トークンの有効期限が切れています"
	case ErrNotFound:
		return "データが見つかりません"
	default:
		return "サーバーエラーが発生しました"
	}
}
