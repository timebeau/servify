package async

import (
	"context"
	"sync"
	"time"
)

// DeadLetterEntry records a failed event for later analysis or replay.
type DeadLetterEntry struct {
	EventID    string    `json:"event_id"`
	EventType  string    `json:"event_type"`
	Error      string    `json:"error"`
	OccurredAt time.Time `json:"occurred_at"`
	Retries    int       `json:"retries"`
}

// DeadLetterRecorder stores dead-lettered events.
type DeadLetterRecorder interface {
	Record(ctx context.Context, entry DeadLetterEntry) error
	List(ctx context.Context, eventType string, limit int) ([]DeadLetterEntry, error)
}

// InMemoryDeadLetterRecorder is an in-memory implementation for development.
type InMemoryDeadLetterRecorder struct {
	mu      sync.Mutex
	entries []DeadLetterEntry
	limit   int
}

// NewInMemoryDeadLetterRecorder creates a new in-memory dead letter recorder.
// The limit parameter controls the maximum number of entries retained.
func NewInMemoryDeadLetterRecorder(limit int) *InMemoryDeadLetterRecorder {
	if limit <= 0 {
		limit = 1000
	}
	return &InMemoryDeadLetterRecorder{
		entries: make([]DeadLetterEntry, 0, limit),
		limit:   limit,
	}
}

// Record appends a dead letter entry, evicting oldest entries if at capacity.
func (r *InMemoryDeadLetterRecorder) Record(_ context.Context, entry DeadLetterEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.entries) >= r.limit {
		r.entries = r.entries[1:]
	}
	r.entries = append(r.entries, entry)
	return nil
}

// List returns dead letter entries, optionally filtered by event type.
func (r *InMemoryDeadLetterRecorder) List(_ context.Context, eventType string, limit int) ([]DeadLetterEntry, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if limit <= 0 || limit > len(r.entries) {
		limit = len(r.entries)
	}

	if eventType == "" {
		result := make([]DeadLetterEntry, limit)
		copy(result, r.entries[:limit])
		return result, nil
	}

	var filtered []DeadLetterEntry
	for _, e := range r.entries {
		if e.EventType == eventType {
			filtered = append(filtered, e)
			if len(filtered) >= limit {
				break
			}
		}
	}
	return filtered, nil
}
