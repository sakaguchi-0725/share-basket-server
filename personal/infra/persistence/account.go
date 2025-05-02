package persistence

import (
	"share-basket-server/personal/domain"
	"share-basket-server/personal/infra/dto"

	"gorm.io/gorm"
)

type accountPersistence struct {
	db *gorm.DB
}

func (a *accountPersistence) Store(acc *domain.Account) error {
	accDto := dto.NewAccountDto(*acc)

	err := a.db.Save(&accDto).Error
	if err != nil {
		return err
	}

	updatedAccount := accDto.ToModel()
	*acc = updatedAccount

	return nil
}

func NewAccountPersistence(db *gorm.DB) domain.AccountRepository {
	return &accountPersistence{db}
}
