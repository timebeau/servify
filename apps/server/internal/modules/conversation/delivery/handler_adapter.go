package delivery

import (
	"context"
	"errors"

	conversationapp "servify/apps/server/internal/modules/conversation/application"
	conversationdomain "servify/apps/server/internal/modules/conversation/domain"

	"gorm.io/gorm"
)

type HandlerServiceAdapter struct {
	service *conversationapp.Service
}

func NewHandlerService(service *conversationapp.Service) *HandlerServiceAdapter {
	return &HandlerServiceAdapter{service: service}
}

func (a *HandlerServiceAdapter) ListMessages(ctx context.Context, sessionID string, limit int) ([]conversationapp.ConversationMessageDTO, error) {
	items, err := a.service.ListRecentMessages(ctx, sessionID, limit)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrConversationNotFound
		}
		return nil, err
	}
	return items, nil
}

func (a *HandlerServiceAdapter) SendAgentMessage(ctx context.Context, sessionID string, content string) (*conversationapp.ConversationMessageDTO, error) {
	if _, err := a.service.GetConversation(ctx, sessionID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrConversationNotFound
		}
		return nil, err
	}

	item, err := a.service.IngestTextMessage(ctx, conversationapp.IngestTextMessageCommand{
		ConversationID: sessionID,
		Sender:         conversationdomain.ParticipantRoleAgent,
		Content:        content,
		Metadata: map[string]string{
			"source": "admin",
		},
	})
	if err != nil {
		return nil, err
	}
	return item, nil
}

var _ HandlerService = (*HandlerServiceAdapter)(nil)
