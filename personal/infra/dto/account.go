package dto

import (
	"share-basket-server/personal/domain"
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

func NewAccountDto(acc domain.Account) Account {
	return Account{
		ID:     acc.ID.String(),
		UserID: acc.UserID.String(),
		Name:   acc.Name,
	}
}

func (acc Account) ToModel() domain.Account {
	return domain.RecreateAccount(
		domain.AccountID(acc.ID),
		domain.UserID(acc.UserID),
		acc.Name,
	)
}
