package persistence

import (
	"share-basket-server/personal/domain/model"
	"share-basket-server/personal/domain/repository"
	"share-basket-server/personal/infra/dto"

	"gorm.io/gorm"
)

type accountPersistence struct {
	db *gorm.DB
}

func (a *accountPersistence) Store(acc *model.Account) error {
	accDto := dto.NewAccountDto(*acc)

	err := a.db.Save(&accDto).Error
	if err != nil {
		return err
	}

	updatedAccount := accDto.ToModel()
	*acc = updatedAccount

	return nil
}

func NewAccountPersistence(db *gorm.DB) repository.Account {
	return &accountPersistence{db}
}
