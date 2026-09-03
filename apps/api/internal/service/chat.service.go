package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/lucassoaresfr/winthor-sql-agent.git/internal/domain"
)

type chatService struct {
	repo domain.ChatRepository
}

func NewChatService(repo domain.ChatRepository) domain.ChatService {
	return &chatService{repo: repo}
}

func (s *chatService) StartNewChat(c context.Context, input domain.ToCreateChatInput) (*domain.Chat, error) {
	if input.Title == "" {
		input.Title = "Nova Conversa"
	}

	return s.repo.ToCreateChat(c, input)
}

func (s *chatService) ToAddMessage(c context.Context, input domain.ToCreateMesageInput) (*domain.Message, error) {
	return s.repo.SaveMessage(c, input)
}

func (s *chatService) GethistoryChat(c context.Context, chatID uuid.UUID) (*domain.Chat, error) {
	return s.repo.SearchChatID(c, chatID)
}

func (s *chatService) ListConversationsUser(c context.Context, userID string) ([]domain.Chat, error) {
	return s.repo.ListChatsUser(c, userID)
}
