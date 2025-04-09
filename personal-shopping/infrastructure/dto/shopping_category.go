package dto

import (
	"share-basket/personal-shopping/domain/model"
	"time"
)

type ShoppingCategoryDto struct {
	ID        *int64    `gorm:"primaryKey autoIncrement"`
	Name      string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func NewShoppingCategoryDto(model model.ShoppingCategory) ShoppingCategoryDto {
	return ShoppingCategoryDto{
		ID:   model.ID,
		Name: model.Name,
	}
}

func (dto ShoppingCategoryDto) ToModel() model.ShoppingCategory {
	return model.NewShoppingCategory(dto.ID, dto.Name)
}
