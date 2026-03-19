package delivery

import "time"

// WebRTCMediaPayload is the normalized input payload accepted by the WebRTC media adapter.
type WebRTCMediaPayload struct {
	CallID         string
	ConversationID string
	ConnectionID   string
	OccurredAt     time.Time
	Metadata       map[string]interface{}
}
