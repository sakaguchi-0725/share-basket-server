package db

import (
	"fmt"
	"sharebasket/core"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type (
	Conn struct {
		*gorm.DB
	}

	config struct {
		port     string
		host     string
		name     string
		user     string
		password string
	}
)

// 新しいデータベース接続を作成。
// 接続に失敗した場合はエラーを返す。
func New() (*Conn, error) {
	cfg := newConfig()

	db, err := gorm.Open(postgres.Open(cfg.dsn()), &gorm.Config{})
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

func newConfig() config {
	return config{
		port:     core.GetEnv("DB_PORT", "5432"),
		host:     core.GetEnv("DB_HOST", "db"),
		name:     core.GetEnv("DB_NAME", "share-basket"),
		user:     core.GetEnv("DB_USER", "postgres"),
		password: core.GetEnv("DB_PASSWORD", "postgres"),
	}
}

func (c config) dsn() string {
	return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		c.host, c.port, c.name, c.user, c.password)
}
