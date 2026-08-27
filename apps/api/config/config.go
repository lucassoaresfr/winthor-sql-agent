package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	APIAuthToken string
	DB           DBConfig
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	port := getEnv("PORT", "8080")
	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "1521"))

	return &Config{
		Port:         port,
		APIAuthToken: getEnv("API_AUTH_TOKEN", ""),
		DB: DBConfig{
			Host:        getEnv("DB_HOST", "localhost"),
			Port:        dbPort,
			User:        getEnv("DB_USER", ""),
			Password:    getEnv("DB_PASSWORD", ""),
			ServiceName: getEnv("DB_SERVICE_NAME", "ORCL"),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return defaultValue
}
