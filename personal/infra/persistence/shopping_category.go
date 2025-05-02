package persistence

import (
	"share-basket-server/personal/domain"
	"share-basket-server/personal/infra/dto"

	"gorm.io/gorm"
)

type shoppingCategory struct {
	db *gorm.DB
}

func (s *shoppingCategory) GetAll() ([]domain.ShoppingCategory, error) {
	var categories dto.ShoppingCategories

	err := s.db.Find(&categories).Error
	if err != nil {
		return []domain.ShoppingCategory{}, err
	}

	return categories.ToModels(), nil
}

func NewShoppingCategory(db *gorm.DB) domain.ShoppingCategoryRepository {
	return &shoppingCategory{db}
}
