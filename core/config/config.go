package config

import "os"

type App struct {
	Env string
	DB  DB
	AWS AWS
}

func Load() App {
	return App{
		Env: getEnv("APP_ENV", "dev"),
		DB:  newDBConfig(),
		AWS: newAWSConfig(),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
