package bootstrap

import (
	"context"
	"log"

	"github.com/go-redis/redis"
	"github.com/lucassoaresfr/winthor-sql-agent.git/client"
	"github.com/lucassoaresfr/winthor-sql-agent.git/config"
	"github.com/lucassoaresfr/winthor-sql-agent.git/register"
	"google.golang.org/genai"
)

type Application struct {
	Config *config.Config
}

func NewApp() *Application {
	cfg := config.LoadConfig()
	if cfg.GeminiAPIKey == "" {
		log.Fatal("GEMINI_API_KEY é obrigatória.")
	}

	return &Application{
		Config: cfg,
	}
}

func (a *Application) Run() {
	ctx := context.Background()

	// 1. Inicializa Cliente Gemini
	geminiClient, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: a.Config.GeminiAPIKey,
	})
	if err != nil {
		log.Fatalf("Falha ao iniciar cliente Gemini: %v", err)
	}

	// 2. Inicializa Cliente HTTP
	apiClient := client.NewAPIClient(a.Config.APIBaseUrl, a.Config.APIAuthToken)

	rdb := redis.NewClient(&redis.Options{
		Addr:     a.Config.RedisAddr,
		Password: a.Config.RedisPassword,
		DB:       0,
	})

	// 3. Registra todas as dependências e rotas
	router := register.RegisterDependencies(geminiClient, apiClient, a.Config, rdb)

	// 4. Sobe o servidor HTTP
	log.Printf("Servidor orquestrador rodando na porta %s...", a.Config.Port)
	if err := router.Run(":" + a.Config.Port); err != nil {
		log.Fatalf("Erro ao rodar servidor HTTP: %v", err)
	}
}
