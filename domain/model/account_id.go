package model

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrInvalidAccountID = errors.New("invalid account id")
)

type AccountID string

func NewAccountID() AccountID {
	return AccountID(uuid.NewString())
}

func ParseAccountID(s string) (AccountID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return "", ErrInvalidAccountID
	}

	return AccountID(id.String()), nil
}

func (a AccountID) String() string {
	return string(a)
}
