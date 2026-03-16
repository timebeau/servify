package application

import (
	"context"
	"fmt"
)

type HandoffPort interface {
	RequestHandoff(ctx context.Context, conversationID string, reason string) (map[string]interface{}, error)
}

type HandoffTool struct {
	port HandoffPort
}

func NewHandoffTool(port HandoffPort) *HandoffTool {
	return &HandoffTool{port: port}
}

func (t *HandoffTool) Name() string { return "handoff" }

func (t *HandoffTool) Description() string {
	return "Request a handoff to a human agent."
}

func (t *HandoffTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"conversation_id": map[string]interface{}{"type": "string"},
			"reason":          map[string]interface{}{"type": "string"},
		},
		"required": []string{"conversation_id"},
	}
}

func (t *HandoffTool) Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	if t.port == nil {
		return nil, fmt.Errorf("handoff port is not configured")
	}
	rawID, ok := input["conversation_id"]
	if !ok {
		return nil, fmt.Errorf("missing conversation_id")
	}
	conversationID, ok := rawID.(string)
	if !ok || conversationID == "" {
		return nil, fmt.Errorf("invalid conversation_id")
	}
	reason, _ := input["reason"].(string)
	return t.port.RequestHandoff(ctx, conversationID, reason)
}
