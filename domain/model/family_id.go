package model

import (
	"fmt"
	"sharebasket/core"

	"github.com/google/uuid"
)

type FamilyID string

func NewFamilyID() FamilyID {
	return FamilyID(uuid.NewString())
}

func ParseFamilyID(s string) (FamilyID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return "", core.NewInvalidError(fmt.Errorf("invalid family ID: %w", err))
	}

	return FamilyID(id.String()), nil
}

func (f FamilyID) String() string {
	return string(f)
}
