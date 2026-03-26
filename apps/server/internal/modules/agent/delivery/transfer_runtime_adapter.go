package delivery

import (
	"context"

	platformauth "servify/apps/server/internal/platform/auth"

	"gorm.io/gorm"
)

type TransferRuntimeAdapter struct{}

func NewTransferRuntimeAdapter() *TransferRuntimeAdapter {
	return &TransferRuntimeAdapter{}
}

func (a *TransferRuntimeAdapter) SyncTransferLoad(ctx context.Context, tx *gorm.DB, fromAgentID *uint, toAgentID uint) error {
	updateAgents := func(query *gorm.DB) *gorm.DB {
		if tenantID := platformauth.TenantIDFromContext(ctx); tenantID != "" {
			query = query.Where("tenant_id = ?", tenantID)
		}
		if workspaceID := platformauth.WorkspaceIDFromContext(ctx); workspaceID != "" {
			query = query.Where("workspace_id = ?", workspaceID)
		}
		return query
	}
	if fromAgentID != nil && *fromAgentID != toAgentID {
		if err := updateAgents(tx.WithContext(ctx).Table("agents")).
			Where("user_id = ?", *fromAgentID).
			UpdateColumn("current_load", gorm.Expr("CASE WHEN current_load > 0 THEN current_load - 1 ELSE 0 END")).Error; err != nil {
			return err
		}
	}
	return updateAgents(tx.WithContext(ctx).Table("agents")).
		Where("user_id = ?", toAgentID).
		UpdateColumn("current_load", gorm.Expr("current_load + 1")).Error
}

var _ RuntimeService = (*TransferRuntimeAdapter)(nil)
