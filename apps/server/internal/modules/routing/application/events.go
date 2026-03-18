package application

import (
	"context"
	"fmt"
	"time"

	"servify/apps/server/internal/platform/eventbus"
)

const (
	RoutingAgentAssignedEventName     = "routing.agent_assigned"
	RoutingTransferCompletedEventName = "routing.transfer_completed"
)

type RoutingEvent struct {
	eventbus.BaseEvent
	SessionID string
	Payload   interface{}
}

func NewRoutingEvent(name string, sessionID string, payload interface{}) RoutingEvent {
	return RoutingEvent{
		BaseEvent: eventbus.BaseEvent{
			EventID:          fmt.Sprintf("%s-%s-%d", name, sessionID, time.Now().UnixNano()),
			EventName:        name,
			EventOccurredAt:  time.Now(),
			EventAggregateID: fmt.Sprintf("routing:%s", sessionID),
		},
		SessionID: sessionID,
		Payload:   payload,
	}
}

type EventPublisher interface {
	Publish(ctx context.Context, event eventbus.Event) error
}
