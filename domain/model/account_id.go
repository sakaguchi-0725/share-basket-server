package model

import (
	"fmt"
	"sharebasket/core"

	"github.com/google/uuid"
)

type AccountID string

func NewAccountID() AccountID {
	return AccountID(uuid.NewString())
}

func ParseAccountID(s string) (AccountID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return "", core.NewInvalidError(fmt.Errorf("invalid account ID: %w", err))
	}

	return AccountID(id.String()), nil
}

func (a AccountID) String() string {
	return string(a)
}
