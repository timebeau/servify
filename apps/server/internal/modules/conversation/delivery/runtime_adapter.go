package delivery

import (
	"context"
	"fmt"
	"time"

	"servify/apps/server/internal/models"
	conversationdomain "servify/apps/server/internal/modules/conversation/domain"
	conversationinfra "servify/apps/server/internal/modules/conversation/infra"
	"servify/apps/server/internal/platform/eventbus"

	"gorm.io/gorm"
)

type RuntimeAdapter struct {
	publisher eventbus.Bus
	db        *gorm.DB
	now       func() time.Time
}

func NewRuntimeAdapter(db *gorm.DB, publisher eventbus.Bus) *RuntimeAdapter {
	return &RuntimeAdapter{publisher: publisher, db: db, now: time.Now}
}

func (a *RuntimeAdapter) LoadTransferSession(ctx context.Context, sessionID string) (*TransferSession, error) {
	if a.db == nil {
		return nil, gorm.ErrInvalidDB
	}
	var model models.Session
	if err := a.db.WithContext(ctx).
		Preload("User").
		First(&model, "id = ?", sessionID).Error; err != nil {
		return nil, err
	}
	return &TransferSession{
		ID:           model.ID,
		CustomerID:   model.UserID,
		AgentID:      model.AgentID,
		TicketID:     model.TicketID,
		Status:       model.Status,
		Platform:     model.Platform,
		UserName:     model.User.Name,
		UserUsername: model.User.Username,
	}, nil
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
