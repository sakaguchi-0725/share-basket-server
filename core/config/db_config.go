package config

import (
	"fmt"
)

type DB struct {
	Port     string
	Host     string
	Name     string
	User     string
	Password string
}

func newDBConfig() DB {
	return DB{
		Port:     getEnv("DB_PORT", "5432"),
		Host:     getEnv("DB_HOST", "db"),
		Name:     getEnv("POSTGRES_DB", "share-basket"),
		User:     getEnv("POSTGRES_USER", "postgres"),
		Password: getEnv("POSTGRES_PASSWORD", "postgres"),
	}
}

func (db DB) DSN() string {
	return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		db.Host, db.Port, db.Name, db.User, db.Password)
}
