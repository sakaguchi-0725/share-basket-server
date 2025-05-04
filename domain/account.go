//go:generate mockgen -destination=../mock/$GOPACKAGE/$GOFILE . AccountRepository
package domain

import "errors"

var ErrAccountNameRequired = errors.New("account name is required")

type (
	Account struct {
		ID     AccountID
		UserID UserID
		Name   string
	}

	AccountRepository interface {
		Store(acc *Account) error
	}
)

func NewAccount(id AccountID, userID UserID, name string) (Account, error) {
	if name == "" {
		return Account{}, ErrAccountNameRequired
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
