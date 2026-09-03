package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/lucassoaresfr/winthor-sql-agent.git/config"
	"github.com/lucassoaresfr/winthor-sql-agent.git/internal/handler"
)

func NewRouter(handlers *Handlers, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	r.Use(handler.CORSMiddleware())
	// Grupo principal da API V1 protegido pelo AuthMiddleware
	api := r.Group("/api/v1")
	api.Use(handler.AuthMiddleware(cfg.APIAuthToken))

	{
		// -------------------------------------------------------------
		// 1. SUBGRUPO: TOOLS / ERP CONSULTAS (Oracle DB)
		// Endpoints que a IA consome para obter contextos e dados
		// -------------------------------------------------------------
		tools := api.Group("/tools")
		{
			tools.GET("/client", handlers.Client.GetClientes)
			tools.GET("/prod", handlers.Produto.GetProdutos)
			tools.GET("/orders", handlers.Orders.ListarPedidos)
			tools.GET("/items", handlers.Orders.ObterPedidoPorNum)
			tools.GET("/promotion", handlers.Promotion.ListarPromocoes)
		}

		// -------------------------------------------------------------
		// 2. SUBGRUPO: HISTÓRICO & SESSÕES DE CHAT (PostgreSQL DB)
		// Endpoints consumidos pelo frontend para gerenciar conversas
		// -------------------------------------------------------------
		chats := api.Group("/chats")
		{
			chats.POST("", handlers.Chat.ToCreateChat)               // Criar novo chat
			chats.GET("/user/:user_id", handlers.Chat.ListChatsUser) // Listar chats do usuário (Sidebar)
			chats.GET("/:id", handlers.Chat.GetChatID)               // Obter histórico completo de um chat
			chats.POST("/:id/messages", handlers.Chat.ToAddMessage)  // Adicionar mensagem ao histórico
		}
	}

	return r
}
