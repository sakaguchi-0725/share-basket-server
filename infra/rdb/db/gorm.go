package db

import (
	"fmt"
	"share-basket-server/core/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(cfg config.DB) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("database initialize failed: %w", err)
	}

	return db, nil
}
