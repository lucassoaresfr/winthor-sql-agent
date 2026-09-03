package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lucassoaresfr/winthor-sql-agent.git/internal/domain"
	"gorm.io/gorm"
)

type chatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) domain.ChatRepository {
	return &chatRepository{db: db}
}

func (r *chatRepository) ToCreateChat(c context.Context, input domain.ToCreateChatInput) (*domain.Chat, error) {
	chat := &domain.Chat{
		UserID: input.UserID,
		Title:  input.Title,
	}

	err := r.db.WithContext(c).Create(chat).Error

	if err != nil {
		return nil, err
	}

	return chat, nil
}

func (r *chatRepository) SaveMessage(ctx context.Context, input domain.ToCreateMesageInput) (*domain.Message, error) {
	message := domain.Message{
		ID:               uuid.New(), // Preenche o UUID para não enviar '00000000-0000...' ao banco
		ChatID:           input.ChatID,
		Role:             input.Role,
		Content:          input.Content,
		ModelUsed:        input.ModelUsed,
		PromptTokens:     input.PromptTokens,
		CompletionTokens: input.CompletionTokens,
		CreatedAt:        time.Now(),
	}

	// 2. Transação para salvar a mensagem e atualizar o updated_at do chat
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Insere a mensagem
		if err := tx.Create(&message).Error; err != nil {
			return err
		}

		// Atualiza a coluna updated_at na tabela chats
		if err := tx.Model(&domain.Chat{}).
			Where("id = ?", input.ChatID).
			Update("updated_at", time.Now()).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (r *chatRepository) SearchChatID(c context.Context, chatID uuid.UUID) (*domain.Chat, error) {
	var chat domain.Chat

	err := r.db.WithContext(c).Where("id = ?", chatID).First(&chat).Error
	if err != nil {
		return nil, err
	}

	var messages []domain.Message

	err = r.db.WithContext(c).
		Where("chat_id = ?", chatID).
		Order("created_at ASC").
		Find(&messages).Error
	if err != nil {
		return nil, err
	}

	chat.Messages = messages
	return &chat, nil
}

func (r *chatRepository) ListChatsUser(c context.Context, userID string) ([]domain.Chat, error) {
	var chats []domain.Chat

	err := r.db.WithContext(c).Where("user_id", userID).Order("updated_at DESC").Find(&chats).Error

	if err != nil {
		return nil, err
	}

	return chats, nil
}
