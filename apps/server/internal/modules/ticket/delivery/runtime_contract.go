package delivery

import (
	"context"

	"gorm.io/gorm"
)

type RuntimeService interface {
	SyncTransferAssignment(ctx context.Context, tx *gorm.DB, ticketID uint, agentID uint, actorID uint) error
}
