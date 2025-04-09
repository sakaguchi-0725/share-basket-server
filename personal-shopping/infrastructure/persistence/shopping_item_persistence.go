package persistence

import (
	"share-basket/personal-shopping/domain/model"
	"share-basket/personal-shopping/domain/repository"

	"gorm.io/gorm"
)

type shoppingItemPersistence struct {
	db *gorm.DB
}

// Get implements repository.ShoppingItemRepository.
func (s *shoppingItemPersistence) Get(status model.ShoppingStatus) ([]model.ShoppingItem, error) {
	panic("unimplemented")
}

// Save implements repository.ShoppingItemRepository.
func (s *shoppingItemPersistence) Save(item *model.ShoppingItem) error {
	panic("unimplemented")
}

func NewShoppingItemPersistence(db *gorm.DB) repository.ShoppingItemRepository {
	return &shoppingItemPersistence{db}
}
