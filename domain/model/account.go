package model

import "errors"

type Account struct {
	ID        AccountID
	Name      string
	UserID    string
	IsPremium bool
}

// 新しいAccountインスタンスを生成する。
// nameとuserIDが空文字列の場合はエラーを返す。
func NewAccount(id AccountID, name, userID string) (Account, error) {
	if name == "" {
		return Account{}, errors.New("account name is required")
	}

	if userID == "" {
		return Account{}, errors.New("user id is required")
	}

	return Account{
		ID:        id,
		Name:      name,
		UserID:    userID,
		IsPremium: false,
	}, nil
}
