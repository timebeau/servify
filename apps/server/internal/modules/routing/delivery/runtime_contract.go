package delivery

import (
	"context"
	"time"

	"servify/apps/server/internal/models"
)

// RuntimeService is the routing contract used by session-transfer runtime glue.
type RuntimeService interface {
	AddToWaitingQueue(ctx context.Context, sessionID string, reason string, targetSkills []string, priority string, notes string) (*models.WaitingRecord, error)
	GetTransferHistory(ctx context.Context, sessionID string) ([]models.TransferRecord, error)
	ListWaitingRecords(ctx context.Context, status string, limit int) ([]models.WaitingRecord, error)
	GetWaitingRecord(ctx context.Context, sessionID string) (*models.WaitingRecord, error)
	CancelWaiting(ctx context.Context, sessionID string, reason string) (*models.WaitingRecord, error)
	MarkWaitingTransferred(ctx context.Context, sessionID string, agentID uint, assignedAt time.Time) (*models.WaitingRecord, error)
}
