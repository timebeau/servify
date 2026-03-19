package channel

import (
	"fmt"
	"time"
)

const WebChannel = "web"

// NewWebMessageEvent maps browser websocket traffic into the normalized inbound contract.
func NewWebMessageEvent(conversationID, actorID, content string) InboundEvent {
	return InboundEvent{
		EventID:        fmt.Sprintf("web-message-%d", time.Now().UnixNano()),
		Channel:        WebChannel,
		ConversationID: conversationID,
		ActorID:        actorID,
		Kind:           EventKindMessage,
		Payload: map[string]interface{}{
			"content": content,
		},
		OccurredAt: time.Now(),
	}
}
