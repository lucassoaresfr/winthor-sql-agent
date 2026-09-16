package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucassoaresfr/winthor-sql-agent.git/internal/domain"
)

type FinicialHandler struct {
	svc domain.FinicialService
}

func NewFinicialHandler(svc domain.FinicialService) *FinicialHandler {
	return &FinicialHandler{svc: svc}
}

func (h *FinicialHandler) GetLaunch(c *gin.Context) {
	var filter domain.FilterFinicial

	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parâmetros inválidos"})
		return
	}

	launch, summary, err := h.svc.ListLaunch(c.Request.Context(), filter)
	if err != nil {
		log.Printf("[GetLaunch]: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno no servidor"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    launch,
		"total":   len(launch),
		"summary": summary,
	})
}
