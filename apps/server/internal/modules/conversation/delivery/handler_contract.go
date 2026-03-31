package delivery

import (
	"context"
	"errors"

	conversationapp "servify/apps/server/internal/modules/conversation/application"
)

var ErrConversationNotFound = errors.New("conversation not found")

type HandlerService interface {
	GetConversation(ctx context.Context, sessionID string) (*conversationapp.ConversationDTO, error)
	ListMessages(ctx context.Context, sessionID string, limit int) ([]conversationapp.ConversationMessageDTO, error)
	ListMessagesBefore(ctx context.Context, sessionID string, beforeMessageID string, limit int) ([]conversationapp.ConversationMessageDTO, error)
	SendAgentMessage(ctx context.Context, sessionID string, content string) (*conversationapp.ConversationMessageDTO, error)
	AssignAgent(ctx context.Context, sessionID string, agentID uint) (*conversationapp.ConversationDTO, error)
	Transfer(ctx context.Context, sessionID string, toAgentID uint) (*conversationapp.ConversationDTO, error)
	Close(ctx context.Context, sessionID string) (*conversationapp.ConversationDTO, error)
}
