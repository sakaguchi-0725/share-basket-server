package model

import (
	"errors"
	"sharebasket/core"
)

type Account struct {
	ID        AccountID
	Name      string
	UserID    string
	IsPremium bool
}

func NewAccount(id AccountID, name, userID string) (Account, error) {
	if name == "" {
		return Account{}, core.NewInvalidError(errors.New("account name is required"))
	}

	if userID == "" {
		return Account{}, core.NewInvalidError(errors.New("user id is required"))
	}

	return Account{
		ID:        id,
		Name:      name,
		UserID:    userID,
		IsPremium: false,
	}, nil
}
