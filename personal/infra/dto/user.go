package dto

import (
	"share-basket-server/personal/domain/model"
	"time"
)

type User struct {
	ID         string    `gorm:"primaryKey"`
	CognitoUID string    `gorm:"not null"`
	Email      string    `gorm:"not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

func NewUserDto(model model.User) User {
	return User{
		ID:         model.ID.String(),
		CognitoUID: model.CognitoUID,
		Email:      model.Email,
	}
}

func (user User) ToModel() model.User {
	id, _ := model.NewUserID(user.ID)

	return model.User{
		ID:         id,
		CognitoUID: user.CognitoUID,
		Email:      user.Email,
	}
}
