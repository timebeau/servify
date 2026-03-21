package delivery

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type RuntimeService interface {
	SyncTransferAssignment(ctx context.Context, tx *gorm.DB, sessionID string, customerID uint, agentID uint) error
	AppendSystemMessage(ctx context.Context, tx *gorm.DB, sessionID string, content string, createdAt time.Time) error
}
