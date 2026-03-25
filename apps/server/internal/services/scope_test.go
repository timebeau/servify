//go:build integration
// +build integration

package services

import (
	"context"

	platformauth "servify/apps/server/internal/platform/auth"
)

func scopedContext(tenantID, workspaceID string) context.Context {
	return platformauth.ContextWithScope(context.Background(), tenantID, workspaceID)
}
