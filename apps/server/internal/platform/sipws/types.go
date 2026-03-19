package sipws

import "time"

type SignalingMessage struct {
	CallID         string                 `json:"call_id"`
	ConversationID string                 `json:"conversation_id,omitempty"`
	ConnectionID   string                 `json:"connection_id,omitempty"`
	From           string                 `json:"from"`
	To             string                 `json:"to"`
	Method         string                 `json:"method"`
	DTMF           string                 `json:"dtmf,omitempty"`
	OccurredAt     time.Time              `json:"occurred_at"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}
