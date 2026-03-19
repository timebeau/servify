package channel

import "context"

// Adapter defines the lifecycle and mapping responsibilities for external channels.
//
// Responsibilities:
// - normalize provider-specific events into InboundEvent
// - deliver OutboundEvent to the remote provider
// - keep provider session/auth concerns inside the adapter boundary
// - avoid leaking provider DTOs into conversation/ticket modules
type Adapter interface {
	Name() string
	Receive(context.Context) (<-chan InboundEvent, error)
	Send(context.Context, OutboundEvent) error
}
