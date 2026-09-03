package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)



type Config struct {
	Port         string
	APIAuthToken string
	DB           DBConfig       // Oracle
	Postgres     PostgresConfig // PostgreSQL (Histórico)
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	port := getEnv("PORT", "8080")
	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "1521"))
	pgPort, _ := strconv.Atoi(getEnv("PG_PORT", "5432"))

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
		Postgres: PostgresConfig{
			Host:     getEnv("PG_HOST", "localhost"),
			Port:     pgPort,
			User:     getEnv("PG_USER", "postgres"),
			Password: getEnv("PG_PASSWORD", ""),
			DBName:   getEnv("PG_DBNAME", "winthor_agent"),
			SSLMode:  getEnv("PG_SSLMODE", "disable"),
		},
	}, nil
}

// GetPostgresDSN monta a string de conexão para o driver do GORM
func (p *PostgresConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		p.Host, p.User, p.Password, p.DBName, p.Port, p.SSLMode,
	)
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return defaultValue
}
