package delivery

import "time"

// WebRTCMediaPayload is the normalized input payload accepted by the WebRTC media adapter.
type WebRTCMediaPayload struct {
	CallID         string                 `json:"call_id"`
	ConversationID string                 `json:"conversation_id,omitempty"`
	ConnectionID   string                 `json:"connection_id,omitempty"`
	OccurredAt     time.Time              `json:"occurred_at"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}
