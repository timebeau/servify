package delivery

import (
	"context"

	"servify/apps/server/internal/models"
	routingcontract "servify/apps/server/internal/modules/routing/contract"
)

// HandlerService is the only session-transfer contract that HTTP handlers should depend on.
type HandlerService interface {
	TransferToHuman(ctx context.Context, req *routingcontract.TransferRequest) (*routingcontract.TransferResult, error)
	TransferToAgent(ctx context.Context, sessionID string, targetAgentID uint, reason string) (*routingcontract.TransferResult, error)
	GetTransferHistory(ctx context.Context, sessionID string) ([]models.TransferRecord, error)
	ListWaitingRecords(ctx context.Context, status string, limit int) ([]models.WaitingRecord, error)
	CancelWaitingRecord(ctx context.Context, sessionID string, operatorID uint, reason string) error
	ProcessWaitingQueue(ctx context.Context) error
	AutoTransferCheck(ctx context.Context, sessionID string, messages []models.Message) bool
}
