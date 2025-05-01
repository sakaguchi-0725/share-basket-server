package persistence

import (
	"share-basket-server/personal/domain/model"
	"share-basket-server/personal/domain/repository"
	"share-basket-server/personal/infra/dto"

	"gorm.io/gorm"
)

type shoppingCategory struct {
	db *gorm.DB
}

func (s *shoppingCategory) GetAll() ([]model.ShoppingCategory, error) {
	var categories dto.ShoppingCategories

	err := s.db.Find(&categories).Error
	if err != nil {
		return []model.ShoppingCategory{}, err
	}

	return categories.ToModels(), nil
}

func NewShoppingCategory(db *gorm.DB) repository.ShoppingCategory {
	return &shoppingCategory{db}
}
