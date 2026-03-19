package pstnprovider

import "time"

type WebhookEvent struct {
	Provider       string                 `json:"provider"`
	EventID        string                 `json:"event_id"`
	EventType      string                 `json:"event_type"`
	CallID         string                 `json:"call_id"`
	ConversationID string                 `json:"conversation_id,omitempty"`
	From           string                 `json:"from"`
	To             string                 `json:"to"`
	DTMF           string                 `json:"dtmf,omitempty"`
	OccurredAt     time.Time              `json:"occurred_at"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}
