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

func buildParticipantID(prefix string, userID uint) string {
	return fmt.Sprintf("%s:%d", prefix, userID)
}

func uintPtr(v uint) *uint {
	return &v
}

var _ RuntimeService = (*RuntimeAdapter)(nil)
