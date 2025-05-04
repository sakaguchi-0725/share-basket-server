package domain

import (
	"errors"
	"share-basket-server/core/apperr"
)

var (
	ErrCognitoUIDRequired = errors.New("cognito uid is required")
	ErrEmailRequired      = errors.New("email is required")
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
		return User{}, ErrCognitoUIDRequired
	}

	if email == "" {
		return User{}, ErrEmailRequired
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
