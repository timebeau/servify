package delivery

import (
	"context"
	"errors"

	conversationapp "servify/apps/server/internal/modules/conversation/application"
)

var ErrConversationNotFound = errors.New("conversation not found")

type HandlerService interface {
	ListMessages(ctx context.Context, sessionID string, limit int) ([]conversationapp.ConversationMessageDTO, error)
	SendAgentMessage(ctx context.Context, sessionID string, content string) (*conversationapp.ConversationMessageDTO, error)
}
