package application

import (
	"context"

	"servify/apps/server/internal/modules/routing/domain"
)

type RoutingRepository interface {
	CreateAssignment(ctx context.Context, assignment *domain.Assignment) error
	CreateQueueEntry(ctx context.Context, entry *domain.QueueEntry) error
	GetQueueEntry(ctx context.Context, sessionID string) (*domain.QueueEntry, error)
	ListQueueEntries(ctx context.Context, status string, limit int) ([]domain.QueueEntry, error)
	UpdateQueueEntry(ctx context.Context, entry *domain.QueueEntry) error
}
