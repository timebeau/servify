package application

import (
	"time"

	"servify/apps/server/internal/modules/conversation/domain"
)

type CreateConversationCommand struct {
	ConversationID string
	CustomerID     *uint
	Subject        string
	Channel        domain.ChannelBinding
	Participants   []domain.Participant
}

type ResumeConversationQuery struct {
	ConversationID string
}

type IngestTextMessageCommand struct {
	ConversationID string
	MessageID      string
	Sender         domain.ParticipantRole
	Content        string
	Metadata       map[string]string
}

type IngestSystemEventCommand struct {
	ConversationID string
	MessageID      string
	Content        string
	Metadata       map[string]string
}

type ConversationDTO struct {
	ID            string                `json:"id"`
	CustomerID    *uint                 `json:"customer_id,omitempty"`
	Status        string                `json:"status"`
	Subject       string                `json:"subject,omitempty"`
	Channel       domain.ChannelBinding `json:"channel"`
	Participants  []domain.Participant  `json:"participants,omitempty"`
	StartedAt     time.Time             `json:"started_at"`
	LastMessageAt *time.Time            `json:"last_message_at,omitempty"`
	EndedAt       *time.Time            `json:"ended_at,omitempty"`
}

type ConversationMessageDTO struct {
	ID             string            `json:"id"`
	ConversationID string            `json:"conversation_id"`
	Sender         string            `json:"sender"`
	Kind           string            `json:"kind"`
	Content        string            `json:"content"`
	Metadata       map[string]string `json:"metadata,omitempty"`
	CreatedAt      time.Time         `json:"created_at"`
}

func MapConversation(item domain.Conversation) ConversationDTO {
	return ConversationDTO{
		ID:            item.ID,
		CustomerID:    item.CustomerID,
		Status:        string(item.Status),
		Subject:       item.Subject,
		Channel:       item.Channel,
		Participants:  item.Participants,
		StartedAt:     item.StartedAt,
		LastMessageAt: item.LastMessageAt,
		EndedAt:       item.EndedAt,
	}
}

func MapMessage(item domain.ConversationMessage) ConversationMessageDTO {
	return ConversationMessageDTO{
		ID:             item.ID,
		ConversationID: item.ConversationID,
		Sender:         string(item.Sender),
		Kind:           string(item.Kind),
		Content:        item.Content,
		Metadata:       item.Metadata,
		CreatedAt:      item.CreatedAt,
	}
}
