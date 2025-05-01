package dto

import (
	"share-basket-server/personal/domain/model"
	"time"
)

type ShoppingCategory struct {
	ID        uint      `gorm:"primaryKey autoIncrement"`
	Name      string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdatedTime"`
}

type ShoppingCategories []ShoppingCategory

func (categories ShoppingCategories) ToModels() []model.ShoppingCategory {
	models := make([]model.ShoppingCategory, len(categories))

	for i, c := range categories {
		models[i] = model.NewShoppingCategory(c.ID, c.Name)
	}

	return models
}
