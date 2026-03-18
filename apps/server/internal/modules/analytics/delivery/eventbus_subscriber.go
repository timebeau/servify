package delivery

import (
	"context"
	"time"

	analyticsapp "servify/apps/server/internal/modules/analytics/application"
	"servify/apps/server/internal/platform/eventbus"
)

type EventBusSubscriber struct {
	service *analyticsapp.Service
}

func NewEventBusSubscriber(service *analyticsapp.Service) *EventBusSubscriber {
	return &EventBusSubscriber{service: service}
}

func (s *EventBusSubscriber) Register(bus eventbus.Bus) {
	if bus == nil || s == nil || s.service == nil {
		return
	}
	registrations := map[string]analyticsapp.IncrementKind{
		"conversation.created":          analyticsapp.IncrementSessions,
		"conversation.message_received": analyticsapp.IncrementMessages,
		"ticket.created":                analyticsapp.IncrementTickets,
		"ticket.closed":                 analyticsapp.IncrementResolved,
		"sla.violation":                 analyticsapp.IncrementSLA,
	}
	for name, kind := range registrations {
		eventName := name
		incrementKind := kind
		bus.Subscribe(eventName, eventbus.HandlerFunc(func(ctx context.Context, event eventbus.Event) error {
			return s.service.IncrementDailyStat(ctx, analyticsapp.IncrementEvent{
				Date: time.Now(),
				Kind: incrementKind,
			})
		}))
	}
}
