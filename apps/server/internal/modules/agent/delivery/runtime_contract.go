package delivery

import (
	"context"

	"gorm.io/gorm"
)

type RuntimeService interface {
	SyncTransferLoad(ctx context.Context, tx *gorm.DB, fromAgentID *uint, toAgentID uint) error
}
