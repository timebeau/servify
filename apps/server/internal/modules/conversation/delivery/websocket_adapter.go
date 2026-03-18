package delivery

import (
	"context"
	"errors"

	"servify/apps/server/internal/models"
	conversationapp "servify/apps/server/internal/modules/conversation/application"
	conversationdomain "servify/apps/server/internal/modules/conversation/domain"

	"gorm.io/gorm"
)

type WebSocketMessageWriter interface {
	PersistTextMessage(ctx context.Context, sessionID string, content string) error
	HasActiveHumanAgent(ctx context.Context, sessionID string) (bool, error)
	ListRecentMessages(ctx context.Context, sessionID string, limit int) ([]models.Message, error)
}

type WebSocketMessageAdapter struct {
	service *conversationapp.Service
}

func NewWebSocketMessageAdapter(service *conversationapp.Service) *WebSocketMessageAdapter {
	return &WebSocketMessageAdapter{service: service}
}

func (a *WebSocketMessageAdapter) PersistTextMessage(ctx context.Context, sessionID string, content string) error {
	if _, err := a.service.ResumeConversation(ctx, conversationapp.ResumeConversationQuery{
		ConversationID: sessionID,
	}); err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if _, createErr := a.service.CreateConversation(ctx, conversationapp.CreateConversationCommand{
			ConversationID: sessionID,
			Channel: conversationdomain.ChannelBinding{
				Channel:   "web",
				SessionID: sessionID,
			},
		}); createErr != nil {
			return createErr
		}
	}

	_, err := a.service.IngestTextMessage(ctx, conversationapp.IngestTextMessageCommand{
		ConversationID: sessionID,
		Sender:         conversationdomain.ParticipantRoleCustomer,
		Content:        content,
		Metadata: map[string]string{
			"source": "websocket",
		},
	})
	return err
}

func (a *WebSocketMessageAdapter) HasActiveHumanAgent(ctx context.Context, sessionID string) (bool, error) {
	conversation, err := a.service.GetConversation(ctx, sessionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	if conversation.Status == string(conversationdomain.ConversationStatusClosed) {
		return false, nil
	}
	for _, participant := range conversation.Participants {
		if participant.Role == conversationdomain.ParticipantRoleAgent && participant.UserID != nil {
			return true, nil
		}
	}
	return false, nil
}

func (a *WebSocketMessageAdapter) ListRecentMessages(ctx context.Context, sessionID string, limit int) ([]models.Message, error) {
	items, err := a.service.ListRecentMessages(ctx, sessionID, limit)
	if err != nil {
		return nil, err
	}
	out := make([]models.Message, 0, len(items))
	for _, item := range items {
		out = append(out, models.Message{
			SessionID: item.ConversationID,
			Content:   item.Content,
			Type:      item.Kind,
			Sender:    item.Sender,
			CreatedAt: item.CreatedAt,
		})
	}
	return out, nil
}
