package eventbus

import "context"

// Handler consumes an Event.
type Handler interface {
	Handle(ctx context.Context, event Event) error
}

// HandlerFunc adapts a function to the Handler interface.
type HandlerFunc func(ctx context.Context, event Event) error

func (f HandlerFunc) Handle(ctx context.Context, event Event) error {
	return f(ctx, event)
}

// Bus is the event bus contract used by modules.
type Bus interface {
	Subscribe(eventName string, handler Handler)
	Publish(ctx context.Context, event Event) error
}
