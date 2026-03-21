package delivery

import (
	"context"

	"gorm.io/gorm"
)

type RuntimeService interface {
	SyncTransferAssignment(ctx context.Context, tx *gorm.DB, sessionID string, customerID uint, agentID uint) error
}
