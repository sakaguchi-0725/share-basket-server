package database

import (
	"errors"
	"share-basket-server/core/apperr"
	"share-basket-server/domain"
	"time"

	"gorm.io/gorm"
)

type (
	userPersistence struct {
		db *gorm.DB
	}

	userDto struct {
		ID         string    `gorm:"primaryKey"`
		CognitoUID string    `gorm:"not null"`
		Email      string    `gorm:"not null"`
		CreatedAt  time.Time `gorm:"autoCreateTime"`
		UpdatedAt  time.Time `gorm:"autoUpdateTime"`
	}
)

func (u *userPersistence) GetByEmail(email string) (domain.User, error) {
	var user userDto

	err := u.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, apperr.ErrDataNotFound
		}
		return domain.User{}, err
	}

	return user.toModel(), nil
}

func (u *userPersistence) Store(user *domain.User) error {
	userDto := newUserDto(*user)

	err := u.db.Save(&userDto).Error
	if err != nil {
		return err
	}

	updatedUser := userDto.toModel()
	*user = updatedUser

	return nil
}

func NewUserPersistence(db *gorm.DB) domain.UserRepository {
	return &userPersistence{db}
}

func newUserDto(user domain.User) userDto {
	return userDto{
		ID:         user.ID.String(),
		CognitoUID: user.CognitoUID,
		Email:      user.Email,
	}
}

func (user userDto) toModel() domain.User {
	return domain.User{
		ID:         domain.UserID(user.ID),
		CognitoUID: user.CognitoUID,
		Email:      user.Email,
	}
}

func (userDto) TableName() string {
	return "users"
}
