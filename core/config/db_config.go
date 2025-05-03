package config

import (
	"fmt"
	"share-basket-server/core/util"
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
		Port:     util.GetEnv("DB_PORT", "5432"),
		Host:     util.GetEnv("DB_HOST", "db"),
		Name:     util.GetEnv("POSTGRES_DB", "share-basket"),
		User:     util.GetEnv("POSTGRES_USER", "postgres"),
		Password: util.GetEnv("POSTGRES_PASSWORD", "postgres"),
	}
}

func (db DB) DSN() string {
	return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		db.Host, db.Port, db.Name, db.User, db.Password)
}
