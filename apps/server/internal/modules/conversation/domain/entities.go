package domain

import "time"

type ConversationStatus string

const (
	ConversationStatusActive       ConversationStatus = "active"
	ConversationStatusWaitingHuman ConversationStatus = "waiting_human"
	ConversationStatusTransferred  ConversationStatus = "transferred"
	ConversationStatusClosed       ConversationStatus = "closed"
)

type ParticipantRole string

const (
	ParticipantRoleCustomer ParticipantRole = "customer"
	ParticipantRoleAgent    ParticipantRole = "agent"
	ParticipantRoleAI       ParticipantRole = "ai"
	ParticipantRoleSystem   ParticipantRole = "system"
)

type MessageKind string

const (
	MessageKindText   MessageKind = "text"
	MessageKindSystem MessageKind = "system"
)

type ChannelBinding struct {
	Channel     string
	ExternalID  string
	SessionID   string
	WorkspaceID string
	Protocol    string
	ProtocolRef string
}

type Participant struct {
	ID          string
	UserID      *uint
	Role        ParticipantRole
	DisplayName string
}

type Conversation struct {
	ID            string
	CustomerID    *uint
	Status        ConversationStatus
	Subject       string
	Channel       ChannelBinding
	Participants  []Participant
	StartedAt     time.Time
	LastMessageAt *time.Time
	EndedAt       *time.Time
}

type ConversationMessage struct {
	ID             string
	ConversationID string
	Sender         ParticipantRole
	Kind           MessageKind
	Content        string
	Metadata       map[string]string
	CreatedAt      time.Time
}
