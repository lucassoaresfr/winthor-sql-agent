package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/lucassoaresfr/winthor-sql-agent.git/config"
	"github.com/lucassoaresfr/winthor-sql-agent.git/internal/handler"
)

func NewRouter(handlers *Handlers, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// Grupo de rotas protegidas pelo Token
	api := r.Group("/api/v1")

	api.Use(handler.AuthMiddleware(cfg.APIAuthToken))

	{
		api.GET("/client", handlers.Client.GetClientes)
		api.GET("/prod", handlers.Produto.GetProdutos)
	}

	return r
}
