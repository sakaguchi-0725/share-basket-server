package dto

import (
	"share-basket-server/personal/domain/model"
	"time"
)

type Account struct {
	ID        string `gorm:"primaryKey"`
	UserID    string `gorm:"not null"`
	User      User   `gorm:"foreignKey:UserID;references:ID"`
	Name      string
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func NewAccountDto(model model.Account) Account {
	return Account{
		ID:     model.ID.String(),
		UserID: model.UserID.String(),
		Name:   model.Name,
	}
}

func (acc Account) ToModel() model.Account {
	id, _ := model.NewAccountID(acc.ID)
	userID, _ := model.NewUserID(acc.UserID)

	return model.RecreateAccount(id, userID, acc.Name)
}
