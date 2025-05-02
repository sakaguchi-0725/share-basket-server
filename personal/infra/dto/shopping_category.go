package dto

import (
	"share-basket-server/personal/domain"
	"time"
)

type ShoppingCategory struct {
	ID        uint      `gorm:"primaryKey autoIncrement"`
	Name      string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdatedTime"`
}

type ShoppingCategories []ShoppingCategory

func (categories ShoppingCategories) ToModels() []domain.ShoppingCategory {
	models := make([]domain.ShoppingCategory, len(categories))

	for i, c := range categories {
		models[i] = domain.NewShoppingCategory(c.ID, c.Name)
	}

	return models
}
