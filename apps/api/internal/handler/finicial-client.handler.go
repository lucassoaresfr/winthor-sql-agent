package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucassoaresfr/winthor-sql-agent.git/internal/domain"
)

type PCPrestHandler struct {
	svc domain.PCPrestService
}

func NewPCPrestHandler(svc domain.PCPrestService) *PCPrestHandler {
	return &PCPrestHandler{svc: svc}
}

func (h *PCPrestHandler) HandleList(c *gin.Context) {
	var filter domain.FilterPCPrest

	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parâmetos invalidos"})
		return
	}

	titles, resumo, err := h.svc.ListPCPrest(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro interno no servidor"})
		log.Printf("[HandleList]: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": titles, "summary": resumo})
}
