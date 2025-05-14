package model

import "errors"

var (
	ErrRequiredUserID = errors.New("user id is required")
	ErrRequiredEmail  = errors.New("email is required")
)

type User struct {
	ID    string
	Email string
}

// 新しいUserインスタンスを作成する。
// idとemail が空文字列の場合はエラーを返す。
func NewUser(id, email string) (User, error) {
	if id == "" {
		return User{}, ErrRequiredUserID
	}

	if email == "" {
		return User{}, ErrRequiredEmail
	}

	return User{
		ID:    id,
		Email: email,
	}, nil
}
