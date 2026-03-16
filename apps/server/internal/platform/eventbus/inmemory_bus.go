package eventbus

import (
	"context"
	"sync"
)

// InMemoryBus is the default in-process event bus implementation.
type InMemoryBus struct {
	mu       sync.RWMutex
	handlers map[string][]Handler
}

// NewInMemoryBus creates a new in-memory event bus.
func NewInMemoryBus() *InMemoryBus {
	return &InMemoryBus{
		handlers: make(map[string][]Handler),
	}
}

func (b *InMemoryBus) Subscribe(eventName string, handler Handler) {
	if eventName == "" || handler == nil {
		return
	}

	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventName] = append(b.handlers[eventName], handler)
}

func (b *InMemoryBus) Publish(ctx context.Context, event Event) error {
	if event == nil {
		return nil
	}

	b.mu.RLock()
	handlers := append([]Handler(nil), b.handlers[event.Name()]...)
	b.mu.RUnlock()

	for _, handler := range handlers {
		if err := handler.Handle(ctx, event); err != nil {
			return err
		}
	}

	return nil
}
