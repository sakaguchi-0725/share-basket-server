package config

import (
	"share-basket-server/core/util"
)

type App struct {
	Env         string
	FrontendURL string
	DB          DB
	AWS         AWS
}

func Load() App {
	return App{
		Env:         util.GetEnv("APP_ENV", "dev"),
		FrontendURL: util.GetEnv("FRONTEND_URL", "http://localhost:5173"),
		DB:          newDBConfig(),
		AWS:         newAWSConfig(),
	}
}
