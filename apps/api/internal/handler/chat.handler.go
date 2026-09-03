package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lucassoaresfr/winthor-sql-agent.git/internal/domain"
)

type ChatHandler struct {
	chatService domain.ChatService
}

func NewChatHandler(chatService domain.ChatService) *ChatHandler {
	return &ChatHandler{chatService: chatService}
}

func (s *ChatHandler) ToCreateChat(c *gin.Context) {
	var req domain.ToCreateChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados invalidos"})
		return
	}

	input := domain.ToCreateChatInput{
		UserID: req.UserID,
		Title:  req.Title,
	}

	chat, err := s.chatService.StartNewChat(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao criar sessão de chat"})
		log.Printf("[ToCreateChat]: %v", err)
		return
	}

	c.JSON(http.StatusCreated, domain.ToChatResponse(*chat))
}

func (s *ChatHandler) ToAddMessage(c *gin.Context) {
	chatIDParam := c.Param("id")
	chatID, err := uuid.Parse(chatIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do chat invalido"})
		return
	}

	var req domain.SaveMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados invalidos"})
		return
	}

	input := domain.ToCreateMesageInput{
		ChatID:           chatID,
		Role:             req.Role,
		Content:          req.Content,
		ModelUsed:        req.ModelUsed,
		PromptTokens:     req.PromptTokens,
		CompletionTokens: req.CompletionTokens,
	}

	msg, err := s.chatService.ToAddMessage(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao salvar menssgem"})
		log.Printf("[ToAddMessage]: %v", err)
		return
	}

	c.JSON(http.StatusCreated, msg)
}

func (s *ChatHandler) GetChatID(c *gin.Context) {
	chatIDParam := c.Param("id")
	chatID, err := uuid.Parse(chatIDParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID chat iválido"})
		return
	}

	chat, err := s.chatService.GethistoryChat(c.Request.Context(), chatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Chat não encontrado"})
		log.Printf("[GetChatID]: %v", err)
		return
	}

	c.JSON(http.StatusOK, domain.ToChatResponse(*chat))
}

func (s *ChatHandler) ListChatsUser(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID é obrigatorio"})
		return
	}

	chats, err := s.chatService.ListConversationsUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao listar chats"})
		log.Printf("[ListChatsUser]: %v", err)
		return
	}

	var response []domain.ChatResponse
	for _, chat := range chats {
		response = append(response, domain.ToChatResponse(chat)) 
	}

	c.JSON(http.StatusOK, response)
}
