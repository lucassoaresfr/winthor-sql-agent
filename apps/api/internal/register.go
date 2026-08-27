package internal

import (
	"github.com/lucassoaresfr/winthor-sql-agent.git/internal/handler"
	"github.com/lucassoaresfr/winthor-sql-agent.git/internal/repository"
	"github.com/lucassoaresfr/winthor-sql-agent.git/internal/service"
	"gorm.io/gorm"
)

type Handlers struct {
	Client  *handler.ClientHandler
	Produto *handler.ProdutosHandler
}

func RegisterDependencies(db *gorm.DB) *Handlers {

	//clientes
	clientRepo := repository.NewClientRepository(db)
	clientService := service.NewClientService(clientRepo)
	clientHandler := handler.NewClientHandler(clientService)

	//produtos
	prodRepo := repository.NewProdRepository(db)
	prodService := service.NewProdService(prodRepo)
	prodHandler := handler.NewProdHandler(prodService)

	return &Handlers{
		Client:  clientHandler,
		Produto: prodHandler,
	}
}
