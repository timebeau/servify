package async

import (
	"context"
	"errors"
)

// ReplayService defines the interface for re-processing dead-lettered events.
// This is a stub for future implementation.
type ReplayService interface {
	// Replay re-publishes a dead-lettered event by its ID.
	Replay(ctx context.Context, eventID string) error

	// ReplayAll re-publishes all dead-lettered events of a given type.
	ReplayAll(ctx context.Context, eventType string) (int, error)
}

// ErrReplayNotImplemented indicates that replay is not yet available.
var ErrReplayNotImplemented = errors.New("replay not implemented")
