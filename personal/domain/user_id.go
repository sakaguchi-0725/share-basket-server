package domain

import (
	"fmt"

	"github.com/google/uuid"
)

type UserID string

func NewUserID() UserID {
	return UserID(uuid.NewString())
}

func ParseUserID(s string) (UserID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return "", fmt.Errorf("invalid user id: %w", err)
	}

	return UserID(id.String()), nil
}

func (id UserID) String() string {
	return string(id)
}
