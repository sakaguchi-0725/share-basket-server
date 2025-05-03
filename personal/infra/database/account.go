package database

import (
	"share-basket-server/personal/domain"
	"time"

	"gorm.io/gorm"
)

type (
	accountPersistence struct {
		db *gorm.DB
	}

	accountDto struct {
		ID        string  `gorm:"primaryKey"`
		UserID    string  `gorm:"not null"`
		User      userDto `gorm:"foreignKey:UserID;references:ID"`
		Name      string
		CreatedAt time.Time `gorm:"autoCreateTime"`
		UpdatedAt time.Time `gorm:"autoUpdateTime"`
	}
)

func (a *accountPersistence) Store(acc *domain.Account) error {
	accDto := newAccountDto(*acc)

	err := a.db.Save(&accDto).Error
	if err != nil {
		return err
	}

	updatedAccount := accDto.toModel()
	*acc = updatedAccount

	return nil
}

func NewAccountPersistence(db *gorm.DB) domain.AccountRepository {
	return &accountPersistence{db}
}

func newAccountDto(acc domain.Account) accountDto {
	return accountDto{
		ID:     acc.ID.String(),
		UserID: acc.UserID.String(),
		Name:   acc.Name,
	}
}

func (acc accountDto) toModel() domain.Account {
	return domain.RecreateAccount(
		domain.AccountID(acc.ID),
		domain.UserID(acc.UserID),
		acc.Name,
	)
}

func (accountDto) TableName() string {
	return "accounts"
}
