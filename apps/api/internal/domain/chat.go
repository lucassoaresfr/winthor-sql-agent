package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Entidades de Domínio
type Chat struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    string    `json:"user_id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Messages  []Message `json:"messages,omitempty"`
}

func (Chat) TableName() string {
	return "chats"
}

type Message struct {
	ID               uuid.UUID `json:"id"`
	ChatID           uuid.UUID `json:"chat_id"`
	Role             string    `json:"role"`
	Content          string    `json:"content"`
	ModelUsed        string    `json:"model_used,omitempty"`
	PromptTokens     int       `json:"prompt_tokens,omitempty"`
	CompletionTokens int       `json:"completion_tokens,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}

func (Message) TableName() string {
	return "messages"
}

// Filtros / DTOs de entrada do domínio
type ToCreateChatInput struct {
	UserID string
	Title  string
}

type ToCreateMesageInput struct {
	ChatID           uuid.UUID
	Role             string
	Content          string
	ModelUsed        string
	PromptTokens     int
	CompletionTokens int
}

// Interfaces de Contrato (Interfaces no Domain)
type ChatRepository interface {
	ToCreateChat(ctx context.Context, input ToCreateChatInput) (*Chat, error)
	SaveMessage(ctx context.Context, input ToCreateMesageInput) (*Message, error)
	SearchChatID(ctx context.Context, chatID uuid.UUID) (*Chat, error)
	ListChatsUser(ctx context.Context, userID string) ([]Chat, error)
}

type ChatService interface {
	StartNewChat(ctx context.Context, input ToCreateChatInput) (*Chat, error)
	ToAddMessage(ctx context.Context, input ToCreateMesageInput) (*Message, error)
	GethistoryChat(ctx context.Context, chatID uuid.UUID) (*Chat, error)
	ListConversationsUser(ctx context.Context, userID string) ([]Chat, error)
}

//DTO HANDLER

type ToCreateChatRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Title  string `json:"title"`
}

type SaveMessageRequest struct {
	Role             string `json:"role" binding:"required,oneof=user assistant system"`
	Content          string `json:"content" binding:"required"`
	ModelUsed        string `json:"model_used"`
	PromptTokens     int    `json:"prompt_tokens"`
	CompletionTokens int    `json:"completion_tokens"`
}

type ChatResponse struct {
	ID        uuid.UUID         `json:"id"`
	UserID    string            `json:"user_id"`
	Title     string            `json:"title"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	Messages  []MessageResponse `json:"messages,omitempty"`
}

type MessageResponse struct {
	ID        uuid.UUID `json:"id"`
	ChatID    uuid.UUID `json:"chat_id"`
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	ModelUsed string    `json:"model_used,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func ToChatResponse(chat Chat) ChatResponse {
	var msgs []MessageResponse
	for _, m := range chat.Messages {
		msgs = append(msgs, MessageResponse{
			ID:        m.ID,
			ChatID:    m.ChatID,
			Role:      m.Role,
			Content:   m.Content,
			ModelUsed: m.ModelUsed,
			CreatedAt: m.CreatedAt,
		})
	}

	// Faz o parse da string para uuid.UUID
	parsedID, _ := uuid.Parse(chat.ID)

	return ChatResponse{
		ID:        parsedID,
		UserID:    chat.UserID,
		Title:     chat.Title,
		CreatedAt: chat.CreatedAt,
		UpdatedAt: chat.UpdatedAt,
		Messages:  msgs,
	}
}
