package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lucassoaresfr/winthor-sql-agent.git/internal/domain"
)

type OrderHandler struct {
	svc domain.OrderService
}

func NewOrderHandler(svc domain.OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

func (h *OrderHandler) ListarPedidos(c *gin.Context) {
	var filtro domain.FiltroOrder
	if err := c.ShouldBindQuery(&filtro); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parâmetros ivalidos"})
		return
	}

	incluirItens := c.Query("incluir_itens") == "true"

	pedidos, err := h.svc.ListarPedidos(c.Request.Context(), filtro, incluirItens)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao consultar pedidos"})
		log.Printf("[ListarPedidos]: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": pedidos, "total": len(pedidos)})
}

func (h *OrderHandler) ObterPedidoPorNum(c *gin.Context) {
	numPedParam := c.Param("numped")
	numPed, err := strconv.Atoi(numPedParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Número de pedido invalido"})
		return
	}

	order, err := h.svc.ObterPedidoPorNum(c.Request.Context(), numPed)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno na api"})
		log.Printf("[ObterPedidoPorNum]: %v", err)
		return
	}

	if order == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pedido não encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": order})
}
