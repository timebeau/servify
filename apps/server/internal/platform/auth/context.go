package auth

import "context"

type contextKey string

const (
	tenantContextKey    contextKey = "servify.auth.tenant_id"
	workspaceContextKey contextKey = "servify.auth.workspace_id"
)

func ContextWithScope(ctx context.Context, tenantID, workspaceID string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if tenantID != "" {
		ctx = context.WithValue(ctx, tenantContextKey, tenantID)
	}
	if workspaceID != "" {
		ctx = context.WithValue(ctx, workspaceContextKey, workspaceID)
	}
	return ctx
}

func TenantIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	v, _ := ctx.Value(tenantContextKey).(string)
	return v
}

func WorkspaceIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	v, _ := ctx.Value(workspaceContextKey).(string)
	return v
}
