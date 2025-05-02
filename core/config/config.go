package config

import "os"

type App struct {
	Env         string
	FrontendURL string
	DB          DB
	AWS         AWS
}

func Load() App {
	return App{
		Env:         getEnv("APP_ENV", "dev"),
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:5173"),
		DB:          newDBConfig(),
		AWS:         newAWSConfig(),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
