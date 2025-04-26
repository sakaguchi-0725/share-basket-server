package model

import "errors"

type Account struct {
	ID     AccountID
	UserID UserID
	Name   string
}

func NewAccount(id AccountID, userID UserID, name string) (Account, error) {
	if name == "" {
		return Account{}, errors.New("name is required")
	}

	return Account{
		ID:     id,
		UserID: userID,
		Name:   name,
	}, nil
}

func RecreateAccount(id AccountID, userID UserID, name string) Account {
	return Account{
		ID:     id,
		UserID: userID,
		Name:   name,
	}
}
