package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucassoaresfr/winthor-sql-agent.git/internal/domain"
)

type ClientHandler struct {
	service domain.ClientService
}

func NewClientHandler(service domain.ClientService) *ClientHandler {
	return &ClientHandler{service: service}
}

func (h *ClientHandler) GetClientes(c *gin.Context) {
	var filtro domain.FiltroCliente

	if err := c.ShouldBindQuery(&filtro); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parâmetros inválidos na busca"})
		return
	}

	clientes, err := h.service.ListarClientes(c.Request.Context(), filtro)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao consultar clientes no banco."})
		log.Printf("[GetClientes]: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": clientes, "total": len(clientes)})
}
