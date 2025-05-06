package repository

import (
	"errors"
	"share-basket-server/core/apperr"
	"share-basket-server/domain"
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

func (a *accountPersistence) FindByUserID(userID domain.UserID) (domain.Account, error) {
	var account accountDto

	err := a.db.Where("user_id = ?", userID.String()).First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Account{}, apperr.ErrDataNotFound
		}
		return domain.Account{}, err
	}

	return account.toModel(), nil
}

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

func newAccountDto(acc domain.Account) accountDto {
	return accountDto{
		ID:     acc.ID.String(),
		UserID: acc.UserID.String(),
		Name:   acc.Name,
	}
}

func (acc accountDto) toModel() domain.Account {
	return domain.Account{
		ID:     domain.AccountID(acc.ID),
		UserID: domain.UserID(acc.UserID),
		Name:   acc.Name,
	}
}

func (accountDto) TableName() string {
	return "accounts"
}

func NewAccountPersistence(db *gorm.DB) domain.AccountRepository {
	return &accountPersistence{db}
}
