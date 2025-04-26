package service

import (
	"errors"
	"share-basket-server/core/apperr"
	"share-basket-server/personal/domain/repository"
)

type UserService interface {
	IsEmailAvailable(email string) (bool, error)
}

type userService struct {
	userRepo repository.User
}

func (service *userService) IsEmailAvailable(email string) (bool, error) {
	_, err := service.userRepo.GetByEmail(email)
	if err != nil {
		if errors.Is(err, apperr.ErrDataNotFound) {
			return true, nil
		}
		return false, err
	}

	return false, nil
}

func NewUserService(userRepo repository.User) UserService {
	return &userService{userRepo}
}
