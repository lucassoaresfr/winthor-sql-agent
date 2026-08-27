package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucassoaresfr/winthor-sql-agent.git/internal/domain"
)

type ProdutosHandler struct {
	svc domain.ProdutoService
}

func NewProdHandler(svc domain.ProdutoService) *ProdutosHandler {
	return &ProdutosHandler{svc: svc}
}

func (h *ProdutosHandler) GetProdutos(c *gin.Context) {
	var filtro domain.FiltroProduto

	if err := c.ShouldBindQuery(&filtro); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parâmetros ivalidos"})
		return
	}

	produtos, err := h.svc.ListarProdutos(c.Request.Context(), filtro)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno no servidor. Veja o log da api"})
		log.Printf("[GetProdutos]: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": produtos, "total": len(produtos)})
}
