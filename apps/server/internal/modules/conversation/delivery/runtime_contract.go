package delivery

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type TransferSession struct {
	ID           string
	CustomerID   uint
	AgentID      *uint
	TicketID     *uint
	Status       string
	Platform     string
	UserName     string
	UserUsername string
}

type RuntimeService interface {
	LoadTransferSession(ctx context.Context, sessionID string) (*TransferSession, error)
	SyncTransferAssignment(ctx context.Context, tx *gorm.DB, sessionID string, customerID uint, agentID uint) error
	SyncWaitingAssignment(ctx context.Context, tx *gorm.DB, sessionID string, customerID uint) error
	AppendSystemMessage(ctx context.Context, tx *gorm.DB, sessionID string, content string, createdAt time.Time) error
}
