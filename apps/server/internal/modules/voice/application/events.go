package application

import (
	"context"
	"time"

	"servify/apps/server/internal/platform/eventbus"
)

const (
	CallStartedEventName = "call.started"
	CallHeldEventName    = "call.held"
	CallResumedEventName = "call.resumed"
	CallTransferredName  = "call.transferred"
	CallEndedEventName   = "call.ended"
)

type Publisher interface {
	Publish(ctx context.Context, event eventbus.Event) error
}

type VoiceEvent struct {
	eventbus.BaseEvent
	Payload interface{}
}

func NewVoiceEvent(name, aggregateID string, payload interface{}) VoiceEvent {
	return VoiceEvent{
		BaseEvent: eventbus.BaseEvent{
			EventID:          aggregateID + ":" + name,
			EventName:        name,
			EventOccurredAt:  time.Now(),
			EventAggregateID: aggregateID,
		},
		Payload: payload,
	}
}
