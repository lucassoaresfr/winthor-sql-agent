package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucassoaresfr/winthor-sql-agent.git/internal/domain"
)

type PromotionHandler struct {
	svc domain.PromotionService
}

func NewPromotionHandler(svc domain.PromotionService) *PromotionHandler {
	return &PromotionHandler{svc: svc}
}

func (h *PromotionHandler) ListarPromocoes(c *gin.Context) {
	var filtro domain.FiltroPromotion

	if err := c.ShouldBindQuery(&filtro); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pârametro invalidos"})
		return
	}

	promotion, err := h.svc.ListarPromocoes(c.Request.Context(), filtro)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao consultar promoções"})
		log.Printf("[ListarPromocoes]: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": promotion, "total": len(promotion)})

}
