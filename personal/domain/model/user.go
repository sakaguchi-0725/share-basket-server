package model

import (
	"errors"
	"fmt"
	"regexp"
)

type User struct {
	ID         UserID
	CognitoUID string
	Email      string
}

func NewUser(id UserID, cognitoUID, email string) (User, error) {
	if cognitoUID == "" {
		return User{}, errors.New("cognito uid is required")
	}

	if err := validateEmail(email); err != nil {
		return User{}, err
	}

	return User{
		ID:         id,
		CognitoUID: cognitoUID,
		Email:      email,
	}, nil
}

func RecreateUser(id UserID, cognitoUID, email string) User {
	return User{
		ID:         id,
		CognitoUID: cognitoUID,
		Email:      email,
	}
}

func validateEmail(email string) error {
	if email == "" {
		return errors.New("email is required")
	}

	pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`

	matched, err := regexp.MatchString(pattern, email)
	if err != nil {
		return fmt.Errorf("failed to validating email: %w", err)
	}

	if !matched {
		return fmt.Errorf("invalid email: %s", email)
	}

	return nil
}
