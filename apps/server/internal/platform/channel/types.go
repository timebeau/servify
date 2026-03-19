package channel

import "time"

type Direction string

const (
	DirectionInbound  Direction = "inbound"
	DirectionOutbound Direction = "outbound"
)

type EventKind string

const (
	EventKindMessage EventKind = "message"
	EventKindTyping  EventKind = "typing"
	EventKindStatus  EventKind = "status"
	EventKindCall    EventKind = "call"
	EventKindSystem  EventKind = "system"
)

// InboundEvent is the normalized server-side representation for external channel input.
type InboundEvent struct {
	EventID        string                 `json:"event_id"`
	Channel        string                 `json:"channel"`
	ConversationID string                 `json:"conversation_id,omitempty"`
	ActorID        string                 `json:"actor_id,omitempty"`
	Kind           EventKind              `json:"kind"`
	Payload        map[string]interface{} `json:"payload,omitempty"`
	OccurredAt     time.Time              `json:"occurred_at"`
}

// OutboundEvent is the normalized server-side representation for messages sent to a channel.
type OutboundEvent struct {
	EventID        string                 `json:"event_id"`
	Channel        string                 `json:"channel"`
	ConversationID string                 `json:"conversation_id,omitempty"`
	TargetID       string                 `json:"target_id,omitempty"`
	Kind           EventKind              `json:"kind"`
	Payload        map[string]interface{} `json:"payload,omitempty"`
	OccurredAt     time.Time              `json:"occurred_at"`
}
