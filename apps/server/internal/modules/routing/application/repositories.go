package application

import (
	"context"
	"time"

	"servify/apps/server/internal/modules/routing/domain"
)

type RoutingRepository interface {
	CreateAssignment(ctx context.Context, assignment *domain.Assignment) error
	ListAssignments(ctx context.Context, sessionID string) ([]domain.TransferRecord, error)
	CreateQueueEntry(ctx context.Context, entry *domain.QueueEntry) error
	GetQueueEntry(ctx context.Context, sessionID string) (*domain.QueueEntry, error)
	ListQueueEntries(ctx context.Context, status string, limit int) ([]domain.QueueEntry, error)
	UpdateQueueEntry(ctx context.Context, entry *domain.QueueEntry) error
	MarkQueueEntryTransferred(ctx context.Context, sessionID string, agentID uint, assignedAt time.Time) (*domain.QueueEntry, error)
}
