package eventbus

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
)

// InMemoryBus is the default in-process event bus implementation.
//
// Runtime boundary: events are NOT persisted. A process restart loses all
// in-flight events. This is an explicit trade-off: consumers that require
// durability should use an out-of-process message queue instead.
type InMemoryBus struct {
	mu       sync.RWMutex
	handlers map[string][]Handler
	logger   *logrus.Logger

	// observability counters
	publishedCount  atomic.Int64
	failedCount     atomic.Int64
	deadLetterCount atomic.Int64

	// deadLetterMu protects the dead letter slice.
	deadLetterMu   sync.Mutex
	deadLetters    []DeadLetterEntry
	maxDeadLetters int
}

// DeadLetterEntry records a handler failure for diagnostics.
type DeadLetterEntry struct {
	EventID    string
	EventName  string
	HandlerIdx int
	Error      string
	FailedAt   time.Time
}

// BusHealth returns a snapshot of bus metrics.
type BusHealth struct {
	PublishedCount  int64
	FailedCount     int64
	DeadLetterCount int64
	DeadLetters     []DeadLetterEntry
}

// NewInMemoryBus creates a new in-memory event bus.
func NewInMemoryBus() *InMemoryBus {
	return &InMemoryBus{
		handlers:       make(map[string][]Handler),
		logger:         logrus.StandardLogger(),
		maxDeadLetters: 100,
	}
}

// NewInMemoryBusWithLogger creates a new in-memory event bus with a custom logger.
func NewInMemoryBusWithLogger(logger *logrus.Logger) *InMemoryBus {
	bus := NewInMemoryBus()
	bus.logger = logger
	return bus
}

func (b *InMemoryBus) Subscribe(eventName string, handler Handler) {
	if eventName == "" || handler == nil {
		return
	}

	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventName] = append(b.handlers[eventName], handler)
}

// Publish dispatches the event to all subscribed handlers.
// Unlike the naive implementation, it continues after a handler error,
// collecting all errors. Failed dispatches are recorded in the dead letter log.
func (b *InMemoryBus) Publish(ctx context.Context, event Event) error {
	if event == nil {
		return nil
	}

	b.publishedCount.Add(1)

	b.mu.RLock()
	handlers := append([]Handler(nil), b.handlers[event.Name()]...)
	b.mu.RUnlock()

	var firstErr error
	for i, handler := range handlers {
		if err := handler.Handle(ctx, event); err != nil {
			b.failedCount.Add(1)
			b.recordDeadLetter(event, i, err)
			b.logger.WithFields(logrus.Fields{
				"event_name":  event.Name(),
				"event_id":    event.ID(),
				"handler_idx": i,
			}).Warnf("event handler failed: %v", err)
			if firstErr == nil {
				firstErr = err
			}
		}
	}

	return firstErr
}

// Health returns current bus metrics and recent dead letter entries.
func (b *InMemoryBus) Health() BusHealth {
	b.deadLetterMu.Lock()
	dl := make([]DeadLetterEntry, len(b.deadLetters))
	copy(dl, b.deadLetters)
	b.deadLetterMu.Unlock()

	return BusHealth{
		PublishedCount:  b.publishedCount.Load(),
		FailedCount:     b.failedCount.Load(),
		DeadLetterCount: b.deadLetterCount.Load(),
		DeadLetters:     dl,
	}
}

func (b *InMemoryBus) recordDeadLetter(event Event, handlerIdx int, err error) {
	entry := DeadLetterEntry{
		EventID:    event.ID(),
		EventName:  event.Name(),
		HandlerIdx: handlerIdx,
		Error:      err.Error(),
		FailedAt:   time.Now(),
	}

	b.deadLetterMu.Lock()
	if len(b.deadLetters) >= b.maxDeadLetters {
		b.deadLetters = b.deadLetters[1:]
	}
	b.deadLetters = append(b.deadLetters, entry)
	b.deadLetterMu.Unlock()
	b.deadLetterCount.Add(1)
}
