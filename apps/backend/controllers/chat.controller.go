package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lucassoaresfr/winthor-sql-agent.git/engine"
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
		log.Printf("[HandleWebhook]: %v", err)

		// 2. Trata queda total dos modelos / alta demanda do Gemini (HTTP 503)
		if errors.Is(err, engine.ErrTodosModelosEmCooldown) {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": engine.ErrTodosModelosEmCooldown.Error(),
			})
			return
		}

		// 3. Outros erros internos inesperados no backend (HTTP 500)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ocorreu um erro interno ao processar sua pergunta. Tente novamente em instantes.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sender":   req.Sender,
		"resposta": resposta,
		"status":   "success",
	})
}
