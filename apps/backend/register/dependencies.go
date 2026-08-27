package register

import (
	"github.com/gin-gonic/gin"
	"github.com/lucassoaresfr/winthor-sql-agent.git/client"
	"github.com/lucassoaresfr/winthor-sql-agent.git/config"
	"github.com/lucassoaresfr/winthor-sql-agent.git/controllers"
	"github.com/lucassoaresfr/winthor-sql-agent.git/engine"
	"github.com/lucassoaresfr/winthor-sql-agent.git/routes"
	"github.com/lucassoaresfr/winthor-sql-agent.git/services"
	"google.golang.org/genai"
)

// RegisterDependencies realiza a fiação dos módulos da aplicação
func RegisterDependencies(geminiClient *genai.Client, apiClient *client.APIClient, cfg *config.Config) *gin.Engine {
	// 0. Instancia o cliente da BrasilAPI
	brasilAPIClient := client.NewBrasilAPIClient()

	// 1. Instancia Engine (passando os 4 parâmetros esperados)
	orchestrator := engine.NewOrchestrator(geminiClient, apiClient, brasilAPIClient, cfg)

	// 2. Instancia Service
	chatService := services.NewChatService(orchestrator)

	// 3. Instancia Controller
	chatController := controllers.NewChatController(chatService)

	// 4. Registra Rotas no Gin e retorna o engine HTTP
	return routes.SetupRouter(chatController, cfg)
}
