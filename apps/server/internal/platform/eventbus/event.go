package eventbus

import "time"

// Event is the minimal domain event contract shared by modules.
type Event interface {
	ID() string
	Name() string
	OccurredAt() time.Time
	TenantID() string
	AggregateID() string
}

// BaseEvent provides common metadata for event implementations.
type BaseEvent struct {
	EventID          string
	EventName        string
	EventOccurredAt  time.Time
	EventTenantID    string
	EventAggregateID string
}

func (e BaseEvent) ID() string            { return e.EventID }
func (e BaseEvent) Name() string          { return e.EventName }
func (e BaseEvent) OccurredAt() time.Time { return e.EventOccurredAt }
func (e BaseEvent) TenantID() string      { return e.EventTenantID }
func (e BaseEvent) AggregateID() string   { return e.EventAggregateID }
