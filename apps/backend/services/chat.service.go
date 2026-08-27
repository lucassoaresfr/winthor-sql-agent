package services

import (
	"context"
	"errors"

	"github.com/lucassoaresfr/winthor-sql-agent.git/engine"
)

type ChatService struct {
	Orchestrator *engine.Orchestrator
}

func NewChatService(orchestrator *engine.Orchestrator) *ChatService {
	return &ChatService{Orchestrator: orchestrator}
}

func (s *ChatService) ProcessarMensagem(c context.Context, messages []engine.ChatMessage) (string, error) {
	if len(messages) == 0 {
		return "", errors.New("o histórico de mensagens não pode estar vazio")
	}

	return s.Orchestrator.ProcessarPergunta(c, messages)
}
