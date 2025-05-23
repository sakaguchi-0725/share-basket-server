package db

import (
	"fmt"
	"sharebasket/core"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Conn struct {
	*gorm.DB
}

// 新しいデータベース接続を作成。
// 接続に失敗した場合はエラーを返す。
func New(cfg core.DBConfig) (*Conn, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("database initialize failed: %w", err)
	}

	// 接続テスト
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("database connection test failed: %w", err)
	}

	return &Conn{DB: db}, nil
}
