package delivery

import (
	"context"

	automationapp "servify/apps/server/internal/modules/automation/application"
	"servify/apps/server/internal/platform/eventbus"
)

type EventBusSubscriber struct {
	service *automationapp.Service
}

func NewEventBusSubscriber(service *automationapp.Service) *EventBusSubscriber {
	return &EventBusSubscriber{service: service}
}

func (s *EventBusSubscriber) Register(bus eventbus.Bus) {
	if bus == nil || s == nil || s.service == nil {
		return
	}
	events := []string{
		"ticket.created",
		"ticket.updated",
		"ticket.assigned",
		"ticket.closed",
		"conversation.created",
		"conversation.message_received",
		"routing.agent_assigned",
		"routing.transfer_completed",
	}
	for _, name := range events {
		eventName := name
		bus.Subscribe(eventName, eventbus.HandlerFunc(func(ctx context.Context, evt eventbus.Event) error {
			s.service.HandleBusEvent(ctx, eventName, evt.AggregateID(), evt)
			return nil
		}))
	}
}
