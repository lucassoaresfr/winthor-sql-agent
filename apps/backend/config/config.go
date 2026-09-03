package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                string
	GeminiAPIKey        string
	APIBaseUrl          string
	APIAuthToken        string
	APIAuthTokenBackend string
	RedisAddr           string
	RedisPassword       string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: Arquivo .env não encontrado. Lendo variáveis do sistema")
	}

	port := os.Getenv("PORT")
	if port != "" {
		port = "8081"
	}

	apiBaseUrl := os.Getenv("API_BASE_URL")
	if apiBaseUrl == "" {
		log.Panic("Erro: url da api não encontrada no .env")
	}

	return &Config{
		Port:                port,
		GeminiAPIKey:        os.Getenv("GEMINI_API_KEY"),
		APIAuthTokenBackend: os.Getenv("API_AUTH_TOKEN_BACKEND"),
		APIBaseUrl:          apiBaseUrl,
		APIAuthToken:        os.Getenv("API_AUTH_TOKEN"),
		RedisAddr:           getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:       os.Getenv("REDIS_PASSWORD"),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}
