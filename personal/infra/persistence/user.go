package persistence

import (
	"errors"
	"share-basket-server/core/apperr"
	"share-basket-server/personal/domain"
	"share-basket-server/personal/infra/dto"

	"gorm.io/gorm"
)

type userPersistence struct {
	db *gorm.DB
}

func (u *userPersistence) GetByEmail(email string) (domain.User, error) {
	var user dto.User

	err := u.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, apperr.ErrDataNotFound
		}
		return domain.User{}, err
	}

	return user.ToModel(), nil
}

func (u *userPersistence) Store(user *domain.User) error {
	userDto := dto.NewUserDto(*user)

	err := u.db.Save(&userDto).Error
	if err != nil {
		return err
	}

	updatedUser := userDto.ToModel()
	*user = updatedUser

	return nil
}

func NewUserPersistence(db *gorm.DB) domain.UserRepository {
	return &userPersistence{db}
}
