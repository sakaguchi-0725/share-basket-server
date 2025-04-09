package dto

import (
	"share-basket/personal-shopping/core/util"
	"share-basket/personal-shopping/domain/model"
	"time"
)

type ShoppingItemDto struct {
	ID         *int64              `gorm:"primaryKey autoIncrement"`
	Name       string              `gorm:"not null"`
	CategoryID int64               `gorm:"not null"`
	Category   ShoppingCategoryDto `gorm:"foreignKey:CategoryID;references:ID"`
	Status     string              `gorm:"not null"`
	CreatedAt  time.Time           `gorm:"autoCreateTime"`
	UpdatedAt  time.Time           `gorm:"autoUpdateTime"`
}

func NewShoppingItemDto(model model.ShoppingItem) ShoppingItemDto {
	return ShoppingItemDto{
		ID:         model.ID,
		Name:       model.Name,
		CategoryID: util.Derefer(model.Category.ID),
		Status:     model.Status.String(),
	}
}

func (dto ShoppingItemDto) ToModel() model.ShoppingItem {
	status, _ := model.NewShoppingStatus(dto.Status)
	category := model.NewShoppingCategory(
		util.Ptr(dto.CategoryID), dto.Category.Name,
	)

	return model.RecreateShoppingItem(
		dto.ID,
		dto.Name,
		category,
		status,
	)
}
