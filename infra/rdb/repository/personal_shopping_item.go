package repository

import (
	"share-basket-server/core/util"
	"share-basket-server/domain"
	"time"

	"gorm.io/gorm"
)

type (
	personalShoppingItemPersistence struct {
		db *gorm.DB
	}

	personalShoppingItemDto struct {
		ID         uint                `gorm:"primaryKey autoIncrement"`
		Name       string              `gorm:"not null"`
		Status     string              `gorm:"not null"`
		AccountID  string              `gorm:"not null"`
		Account    accountDto          `gorm:"foreignKey:AccountID;references:ID"`
		CategoryID uint                `gorm:"not null"`
		Category   shoppingCategoryDto `gorm:"foreignKey:CategoryID;references:ID"`
		CreatedAt  time.Time           `gorm:"autoCreateTime"`
		UpdatedAt  time.Time           `gorm:"autoUpdatedTime"`
	}

	personalShoppingItemDtos []personalShoppingItemDto
)

func (p *personalShoppingItemPersistence) GetAll(accID domain.AccountID, status *domain.ShoppingStatus) ([]domain.PersonalShoppingItem, error) {
	query := p.db.Where("account_id = ?", accID.String())

	if status != nil {
		query = query.Where("status = ?", util.Derefer(status))
	}

	var items personalShoppingItemDtos

	err := query.Find(&items).Error
	if err != nil {
		return []domain.PersonalShoppingItem{}, err
	}

	return items.toModel(), nil
}

func (p *personalShoppingItemPersistence) Store(item *domain.PersonalShoppingItem) error {
	panic("unimplemented")
}

func (dtos personalShoppingItemDtos) toModel() []domain.PersonalShoppingItem {
	items := make([]domain.PersonalShoppingItem, len(dtos))

	for i, v := range dtos {
		items[i] = domain.PersonalShoppingItem{
			ID:         util.Ptr(v.ID),
			Name:       v.Name,
			Status:     domain.ShoppingStatus(v.Status),
			CategoryID: v.CategoryID,
			AccountID:  domain.AccountID(v.AccountID),
		}
	}

	return items
}

func (dto personalShoppingItemDto) TableName() string {
	return "personal_shopping_items"
}

func NewPersonalShoppingItemPersistence(db *gorm.DB) domain.PersonalShoppingItemRepository {
	return &personalShoppingItemPersistence{db}
}
