package application

import (
	"context"

	"servify/apps/server/internal/modules/conversation/domain"
)

type ConversationRepository interface {
	CreateConversation(ctx context.Context, conversation *domain.Conversation) error
	GetConversation(ctx context.Context, conversationID string) (*domain.Conversation, error)
	UpdateConversation(ctx context.Context, conversation *domain.Conversation) error
	AppendMessage(ctx context.Context, message *domain.ConversationMessage) error
	ListRecentMessages(ctx context.Context, conversationID string, limit int) ([]domain.ConversationMessage, error)
	ListMessagesBefore(ctx context.Context, conversationID string, beforeMessageID string, limit int) ([]domain.ConversationMessage, error)
}
