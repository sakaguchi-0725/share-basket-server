package persistence

import (
	"errors"
	"share-basket-server/core/apperr"
	"share-basket-server/personal/domain/model"
	"share-basket-server/personal/domain/repository"
	"share-basket-server/personal/infra/dto"

	"gorm.io/gorm"
)

type userPersistence struct {
	db *gorm.DB
}

func (u *userPersistence) GetByEmail(email string) (model.User, error) {
	var user dto.User

	err := u.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.User{}, apperr.ErrDataNotFound
		}
		return model.User{}, err
	}

	return user.ToModel(), nil
}

func (u *userPersistence) Store(user *model.User) error {
	userDto := dto.NewUserDto(*user)

	err := u.db.Save(&userDto).Error
	if err != nil {
		return err
	}

	updatedUser := userDto.ToModel()
	*user = updatedUser

	return nil
}

func NewUserPersistence(db *gorm.DB) repository.User {
	return &userPersistence{db}
}
