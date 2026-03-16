package application

import (
	"context"
	"fmt"
	"strconv"
)

type TicketLookupPort interface {
	GetTicketSummary(ctx context.Context, ticketID uint) (map[string]interface{}, error)
}

type TicketLookupTool struct {
	port TicketLookupPort
}

func NewTicketLookupTool(port TicketLookupPort) *TicketLookupTool {
	return &TicketLookupTool{port: port}
}

func (t *TicketLookupTool) Name() string { return "ticket_lookup" }

func (t *TicketLookupTool) Description() string {
	return "Lookup a ticket summary by ticket id."
}

func (t *TicketLookupTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"ticket_id": map[string]interface{}{
				"type": "integer",
			},
		},
		"required": []string{"ticket_id"},
	}
}

func (t *TicketLookupTool) Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error) {
	if t.port == nil {
		return nil, fmt.Errorf("ticket lookup port is not configured")
	}
	id, err := getUintInput(input, "ticket_id")
	if err != nil {
		return nil, err
	}
	return t.port.GetTicketSummary(ctx, id)
}

func getUintInput(input map[string]interface{}, key string) (uint, error) {
	raw, ok := input[key]
	if !ok {
		return 0, fmt.Errorf("missing %s", key)
	}
	switch v := raw.(type) {
	case int:
		return uint(v), nil
	case int32:
		return uint(v), nil
	case int64:
		return uint(v), nil
	case float64:
		return uint(v), nil
	case string:
		n, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid %s", key)
		}
		return uint(n), nil
	default:
		return 0, fmt.Errorf("invalid %s", key)
	}
}
