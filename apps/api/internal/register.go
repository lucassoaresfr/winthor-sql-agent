package internal

import (
	"github.com/lucassoaresfr/winthor-sql-agent.git/internal/handler"
	"github.com/lucassoaresfr/winthor-sql-agent.git/internal/repository"
	"github.com/lucassoaresfr/winthor-sql-agent.git/internal/service"
	"gorm.io/gorm"
)

type Handlers struct {
	Client    *handler.ClientHandler
	Produto   *handler.ProdutosHandler
	Orders    *handler.OrderHandler
	Promotion *handler.PromotionHandler
	Chat      *handler.ChatHandler // Adicionado o ChatHandler
}

// Passamos o dbOracle (*gorm.DB) e o dbPostgres (*gorm.DB)
func RegisterDependencies(dbOracle *gorm.DB, dbPostgres *gorm.DB) *Handlers {

	// clientes (Oracle)
	clientRepo := repository.NewClientRepository(dbOracle)
	clientService := service.NewClientService(clientRepo)
	clientHandler := handler.NewClientHandler(clientService)

	// produtos (Oracle)
	prodRepo := repository.NewProdRepository(dbOracle)
	prodService := service.NewProdService(prodRepo)
	prodHandler := handler.NewProdHandler(prodService)

	// pedidos (Oracle)
	orderRepo := repository.NewOrdersRepository(dbOracle)
	orderSvc := service.NewOrderService(orderRepo)
	orderHandler := handler.NewOrderHandler(orderSvc)

	// promoções (Oracle)
	promotionRepo := repository.NewPromotionRepository(dbOracle)
	promotionSvc := service.NewPromotionService(promotionRepo)
	promotionHandler := handler.NewPromotionHandler(promotionSvc)

	// chat-historico (PostgreSQL)
	chatRepo := repository.NewChatRepository(dbPostgres)
	chatService := service.NewChatService(chatRepo)
	chatHandler := handler.NewChatHandler(chatService)

	return &Handlers{
		Client:    clientHandler,
		Produto:   prodHandler,
		Orders:    orderHandler,
		Promotion: promotionHandler,
		Chat:      chatHandler,
	}
}
