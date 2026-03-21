package delivery

import (
	"context"
	"fmt"
	"time"

	conversationdomain "servify/apps/server/internal/modules/conversation/domain"
	conversationinfra "servify/apps/server/internal/modules/conversation/infra"
	"servify/apps/server/internal/platform/eventbus"

	"gorm.io/gorm"
)

type RuntimeAdapter struct {
	publisher eventbus.Bus
	now       func() time.Time
}

func NewRuntimeAdapter(publisher eventbus.Bus) *RuntimeAdapter {
	return &RuntimeAdapter{publisher: publisher, now: time.Now}
}

func (a *RuntimeAdapter) SyncTransferAssignment(ctx context.Context, tx *gorm.DB, sessionID string, customerID uint, agentID uint) error {
	now := a.now()
	item := &conversationdomain.Conversation{
		ID:         sessionID,
		CustomerID: uintPtr(customerID),
		Status:     conversationdomain.ConversationStatusActive,
		Channel: conversationdomain.ChannelBinding{
			Channel:   "web",
			SessionID: sessionID,
		},
		Participants: []conversationdomain.Participant{
			{
				ID:     buildParticipantID("customer", customerID),
				UserID: uintPtr(customerID),
				Role:   conversationdomain.ParticipantRoleCustomer,
			},
			{
				ID:     buildParticipantID("agent", agentID),
				UserID: uintPtr(agentID),
				Role:   conversationdomain.ParticipantRoleAgent,
			},
		},
		StartedAt:     now,
		LastMessageAt: &now,
	}
	return conversationinfra.NewGormRepository(tx).UpdateConversation(ctx, item)
}

func (a *RuntimeAdapter) AppendSystemMessage(ctx context.Context, tx *gorm.DB, sessionID string, content string, createdAt time.Time) error {
	repo := conversationinfra.NewGormRepository(tx)
	conversation, err := repo.GetConversation(ctx, sessionID)
	if err != nil {
		return err
	}

	if createdAt.IsZero() {
		createdAt = a.now()
	}

	message := &conversationdomain.ConversationMessage{
		ID:             fmt.Sprintf("%s-system-%d", sessionID, createdAt.UnixNano()),
		ConversationID: sessionID,
		Sender:         conversationdomain.ParticipantRoleSystem,
		Kind:           conversationdomain.MessageKindSystem,
		Content:        content,
		CreatedAt:      createdAt,
	}
	if err := repo.AppendMessage(ctx, message); err != nil {
		return err
	}

	item := &conversationdomain.Conversation{
		ID:            conversation.ID,
		CustomerID:    conversation.CustomerID,
		Status:        conversationdomain.ConversationStatusActive,
		Subject:       conversation.Subject,
		Channel:       conversation.Channel,
		Participants:  conversation.Participants,
		StartedAt:     conversation.StartedAt,
		LastMessageAt: &createdAt,
		EndedAt:       conversation.EndedAt,
	}
	return repo.UpdateConversation(ctx, item)
}

func buildParticipantID(prefix string, userID uint) string {
	return fmt.Sprintf("%s:%d", prefix, userID)
}

func uintPtr(v uint) *uint {
	return &v
}

var _ RuntimeService = (*RuntimeAdapter)(nil)
