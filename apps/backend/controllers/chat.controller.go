package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lucassoaresfr/winthor-sql-agent.git/schema"
	"github.com/lucassoaresfr/winthor-sql-agent.git/services"
)

type ChatController struct {
	Service *services.ChatService
}

func NewChatController(service *services.ChatService) *ChatController {
	return &ChatController{Service: service}
}

func (ctrl *ChatController) HandleWebhook(c *gin.Context) {
	var req schema.ChatWebhookRequest

	// 1. Faz o bind do JSON e valida se houve erro no parse ou se a lista veio vazia
	if err := c.ShouldBindJSON(&req); err != nil || len(req.Messages) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payload inválido. O campo 'messages' é obrigatório e não pode estar vazio"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 60*time.Second)
	defer cancel()

	resposta, err := ctrl.Service.ProcessarMensagem(ctx, req.Messages)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno no backend. Veja o log"})
		log.Printf("[HandleWebhook]: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sender":   req.Sender,
		"resposta": resposta,
		"status":   "success",
	})
}
