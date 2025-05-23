package core

import (
	"fmt"
	"os"
	"strings"
)

type (
	Config struct {
		Port string
		Env  string

		RedisHost string
		DB        DBConfig
		AWS       AWSConfig
	}

	DBConfig struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}

	AWSConfig struct {
		Region              string
		CognitoUserPoolID   string
		CognitoClientID     string
		CognitoClientSecret string
		AccessKeyID         string
		SecretAccessKey     string
	}
)

// DSN はデータベース接続文字列を返します
func (c DBConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		c.Host, c.Port, c.Name, c.User, c.Password)
}

// 環境変数からConfigを読み込む
// 必須の環境変数が設定されていない場合はエラーを返す
func LoadConfig() (*Config, error) {
	config := &Config{
		Port: getEnvOrDefault("PORT", "8080"),
		Env:  getEnvOrDefault("APP_ENV", "development"),

		RedisHost: mustGetEnv("REDIS_HOST"),

		DB: DBConfig{
			Host:     mustGetEnv("DB_HOST"),
			Port:     getEnvOrDefault("DB_PORT", "5432"),
			User:     mustGetEnv("DB_USER"),
			Password: mustGetEnv("DB_PASSWORD"),
			Name:     mustGetEnv("DB_NAME"),
		},

		AWS: AWSConfig{
			Region:              mustGetEnv("AWS_REGION"),
			CognitoUserPoolID:   mustGetEnv("COGNITO_USER_POOL_ID"),
			CognitoClientID:     mustGetEnv("COGNITO_CLIENT_ID"),
			CognitoClientSecret: mustGetEnv("COGNITO_CLIENT_SECRET"),
			AccessKeyID:         mustGetEnv("AWS_ACCESS_KEY_ID"),
			SecretAccessKey:     mustGetEnv("AWS_SECRET_ACCESS_KEY"),
		},
	}

	var missingEnvVars []string
	for _, err := range validateConfig(config) {
		if err != nil {
			missingEnvVars = append(missingEnvVars, err.Error())
		}
	}

	if len(missingEnvVars) > 0 {
		return nil, fmt.Errorf("required environment variables are not set: %s", strings.Join(missingEnvVars, ", "))
	}

	return config, nil
}

// validateConfig はConfigの値を検証
func validateConfig(config *Config) []error {
	var errors []error

	// 必須の環境変数を検証
	if config.DB.Host == "" {
		errors = append(errors, fmt.Errorf("DB_HOST"))
	}
	if config.DB.User == "" {
		errors = append(errors, fmt.Errorf("DB_USER"))
	}
	if config.DB.Password == "" {
		errors = append(errors, fmt.Errorf("DB_PASSWORD"))
	}
	if config.DB.Name == "" {
		errors = append(errors, fmt.Errorf("DB_NAME"))
	}
	if config.AWS.Region == "" {
		errors = append(errors, fmt.Errorf("AWS_REGION"))
	}
	if config.AWS.CognitoUserPoolID == "" {
		errors = append(errors, fmt.Errorf("COGNITO_USER_POOL_ID"))
	}
	if config.AWS.CognitoClientID == "" {
		errors = append(errors, fmt.Errorf("COGNITO_CLIENT_ID"))
	}
	if config.AWS.CognitoClientSecret == "" {
		errors = append(errors, fmt.Errorf("COGNITO_CLIENT_SECRET"))
	}
	if config.AWS.AccessKeyID == "" {
		errors = append(errors, fmt.Errorf("AWS_ACCESS_KEY_ID"))
	}
	if config.AWS.SecretAccessKey == "" {
		errors = append(errors, fmt.Errorf("AWS_SECRET_ACCESS_KEY"))
	}
	if config.RedisHost == "" {
		errors = append(errors, fmt.Errorf("REDIS_HOST"))
	}

	return errors
}

// 環境変数を取得し、設定されていない場合はデフォルト値を返す
func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// 環境変数を取得する
// 環境変数が設定されていない場合は空文字を返す
func mustGetEnv(key string) string {
	return os.Getenv(key)
}
