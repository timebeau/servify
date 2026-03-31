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

func (a *HandlerServiceAdapter) GetConversation(ctx context.Context, sessionID string) (*conversationapp.ConversationDTO, error) {
	dto, err := a.service.GetConversation(ctx, sessionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrConversationNotFound
		}
		return nil, err
	}
	return dto, nil
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

func (a *HandlerServiceAdapter) ListMessagesBefore(ctx context.Context, sessionID string, beforeMessageID string, limit int) ([]conversationapp.ConversationMessageDTO, error) {
	items, err := a.service.ListMessagesBefore(ctx, sessionID, beforeMessageID, limit)
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

func (a *HandlerServiceAdapter) AssignAgent(ctx context.Context, sessionID string, agentID uint) (*conversationapp.ConversationDTO, error) {
	dto, err := a.service.AssignAgent(ctx, sessionID, agentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrConversationNotFound
		}
		return nil, err
	}
	return dto, nil
}

func (a *HandlerServiceAdapter) Transfer(ctx context.Context, sessionID string, toAgentID uint) (*conversationapp.ConversationDTO, error) {
	dto, err := a.service.Transfer(ctx, sessionID, toAgentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrConversationNotFound
		}
		return nil, err
	}
	return dto, nil
}

func (a *HandlerServiceAdapter) Close(ctx context.Context, sessionID string) (*conversationapp.ConversationDTO, error) {
	dto, err := a.service.Close(ctx, sessionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrConversationNotFound
		}
		return nil, err
	}
	return dto, nil
}

var _ HandlerService = (*HandlerServiceAdapter)(nil)
