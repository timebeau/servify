package services

import (
	"context"

	platformauth "servify/apps/server/internal/platform/auth"

	"gorm.io/gorm"
)

func tenantAndWorkspace(ctx context.Context) (string, string) {
	return platformauth.TenantIDFromContext(ctx), platformauth.WorkspaceIDFromContext(ctx)
}

func applyScopeFilter(tx *gorm.DB, ctx context.Context) *gorm.DB {
	tenantID, workspaceID := tenantAndWorkspace(ctx)
	if tenantID != "" {
		tx = tx.Where("tenant_id = ?", tenantID)
	}
	if workspaceID != "" {
		tx = tx.Where("workspace_id = ?", workspaceID)
	}
	return tx
}
