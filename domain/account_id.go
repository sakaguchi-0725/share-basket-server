package domain

import (
	"fmt"

	"github.com/google/uuid"
)

type AccountID string

func NewAccountID() AccountID {
	return AccountID(uuid.NewString())
}

func ParseAccountID(s string) (AccountID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return "", fmt.Errorf("invalid account id: %w", err)
	}

	return AccountID(id.String()), nil
}

func (id AccountID) String() string {
	return string(id)
}
