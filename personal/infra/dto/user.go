package dto

import (
	"share-basket-server/personal/domain"
	"time"
)

type User struct {
	ID         string    `gorm:"primaryKey"`
	CognitoUID string    `gorm:"not null"`
	Email      string    `gorm:"not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

func NewUserDto(user domain.User) User {
	return User{
		ID:         user.ID.String(),
		CognitoUID: user.CognitoUID,
		Email:      user.Email,
	}
}

func (user User) ToModel() domain.User {
	return domain.User{
		ID:         domain.UserID(user.ID),
		CognitoUID: user.CognitoUID,
		Email:      user.Email,
	}
}
