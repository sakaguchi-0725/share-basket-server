package model

import "errors"

var (
	ErrRequiredAccountName = errors.New("account name is required")
)

type Account struct {
	ID     AccountID
	Name   string
	UserID string
}

// 新しいAccountインスタンスを生成する。
// nameとuserIDが空文字列の場合はエラーを返す。
func NewAccount(id AccountID, name, userID string) (Account, error) {
	if name == "" {
		return Account{}, ErrRequiredAccountName
	}

	if userID == "" {
		return Account{}, ErrRequiredUserID
	}

	return Account{
		ID:     id,
		Name:   name,
		UserID: userID,
	}, nil
}
