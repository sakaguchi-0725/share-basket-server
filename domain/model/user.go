package model

import (
	"errors"
	"sharebasket/core"
)

type User struct {
	ID    string
	Email string
}

// 新しいUserインスタンスを作成する。
// idとemail が空文字列の場合はエラーを返す。
func NewUser(id, email string) (User, error) {
	if id == "" {
		return User{}, core.NewInvalidError(errors.New("user id is required"))
	}

	if email == "" {
		return User{}, core.NewInvalidError(errors.New("email is required"))
	}

	return User{
		ID:    id,
		Email: email,
	}, nil
}
