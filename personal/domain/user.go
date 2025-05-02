package domain

import (
	"errors"
	"fmt"
	"regexp"
	"share-basket-server/core/apperr"
)

type (
	User struct {
		ID         UserID
		CognitoUID string
		Email      string
	}

	UserRepository interface {
		GetByEmail(email string) (User, error)
		Store(user *User) error
	}

	UserService interface {
		IsEmailAvailable(email string) (bool, error)
	}
)

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

type userService struct {
	repo UserRepository
}

// Emailが使用可能か判定する
func (service *userService) IsEmailAvailable(email string) (bool, error) {
	_, err := service.repo.GetByEmail(email)
	if err != nil {
		if errors.Is(err, apperr.ErrDataNotFound) {
			return true, nil
		}
		return false, err
	}

	return false, nil
}

func NewUserService(repo UserRepository) UserService {
	return &userService{repo}
}
