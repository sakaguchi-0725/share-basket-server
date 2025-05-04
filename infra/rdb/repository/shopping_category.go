package repository

import (
	"share-basket-server/domain"
	"time"

	"gorm.io/gorm"
)

type (
	shoppingCategoryPersistence struct {
		db *gorm.DB
	}

	shoppingCategoryDto struct {
		ID        uint      `gorm:"primaryKey autoIncrement"`
		Name      string    `gorm:"not null"`
		CreatedAt time.Time `gorm:"autoCreateTime"`
		UpdatedAt time.Time `gorm:"autoUpdatedTime"`
	}

	shoppingCategoryDtos []shoppingCategoryDto
)

func (s *shoppingCategoryPersistence) GetAll() ([]domain.ShoppingCategory, error) {
	var categories shoppingCategoryDtos

	err := s.db.Find(&categories).Error
	if err != nil {
		return []domain.ShoppingCategory{}, err
	}

	return categories.toModels(), nil
}

func NewShoppingCategoryPersistence(db *gorm.DB) domain.ShoppingCategoryRepository {
	return &shoppingCategoryPersistence{db}
}

func (categories shoppingCategoryDtos) toModels() []domain.ShoppingCategory {
	models := make([]domain.ShoppingCategory, len(categories))

	for i, c := range categories {
		models[i] = domain.NewShoppingCategory(c.ID, c.Name)
	}

	return models
}

func (dto shoppingCategoryDto) TableName() string {
	return "shopping_categories"
}
