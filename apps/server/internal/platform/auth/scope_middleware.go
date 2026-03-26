package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// EnforceRequestScope blocks callers from widening or contradicting tenant/workspace scope.
// Scoped agent/end-user tokens cannot provide ad-hoc tenant/workspace selectors.
// Admin/service callers may provide request scope only when the token itself is not already scoped.
func EnforceRequestScope() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestTenantID, tenantConflict := firstScopeValue(c, "X-Tenant-ID", "tenant_id", "tenant")
		if tenantConflict {
			abortJSON(c, http.StatusBadRequest, "BadRequest", "conflicting tenant scope")
			return
		}
		requestWorkspaceID, workspaceConflict := firstScopeValue(c, "X-Workspace-ID", "workspace_id", "workspace")
		if workspaceConflict {
			abortJSON(c, http.StatusBadRequest, "BadRequest", "conflicting workspace scope")
			return
		}

		claimsTenantID, _ := c.Get("tenant_id")
		claimsWorkspaceID, _ := c.Get("workspace_id")
		authTenantID, _ := claimsTenantID.(string)
		authWorkspaceID, _ := claimsWorkspaceID.(string)

		switch getPrincipalKind(c) {
		case PrincipalAdmin, PrincipalService:
			if !scopeCompatible(authTenantID, requestTenantID) || !scopeCompatible(authWorkspaceID, requestWorkspaceID) {
				abortJSON(c, http.StatusForbidden, "Forbidden", "requested scope does not match token scope")
				return
			}
			projectRequestScope(c, authTenantID, authWorkspaceID, requestTenantID, requestWorkspaceID)
		default:
			if !scopeCompatible(authTenantID, requestTenantID) || !scopeCompatible(authWorkspaceID, requestWorkspaceID) {
				abortJSON(c, http.StatusForbidden, "Forbidden", "requested scope does not match token scope")
				return
			}
			if (authTenantID == "" && requestTenantID != "") || (authWorkspaceID == "" && requestWorkspaceID != "") {
				abortJSON(c, http.StatusForbidden, "Forbidden", "principal is not allowed to widen request scope")
				return
			}
		}

		c.Next()
	}
}

func firstScopeValue(c *gin.Context, headerKey string, queryKeys ...string) (string, bool) {
	var values []string
	if v := strings.TrimSpace(c.GetHeader(headerKey)); v != "" {
		values = append(values, v)
	}
	for _, key := range queryKeys {
		if v := strings.TrimSpace(c.Query(key)); v != "" {
			values = append(values, v)
		}
	}
	if len(values) == 0 {
		return "", false
	}
	base := values[0]
	for _, item := range values[1:] {
		if item != base {
			return "", true
		}
	}
	return base, false
}

func scopeCompatible(authValue, requestedValue string) bool {
	return requestedValue == "" || authValue == "" || authValue == requestedValue
}

func projectRequestScope(c *gin.Context, authTenantID, authWorkspaceID, requestTenantID, requestWorkspaceID string) {
	resolvedTenantID := authTenantID
	if resolvedTenantID == "" {
		resolvedTenantID = requestTenantID
	}
	resolvedWorkspaceID := authWorkspaceID
	if resolvedWorkspaceID == "" {
		resolvedWorkspaceID = requestWorkspaceID
	}
	if resolvedTenantID == "" && resolvedWorkspaceID == "" {
		return
	}
	reqCtx := ContextWithScope(c.Request.Context(), resolvedTenantID, resolvedWorkspaceID)
	c.Request = c.Request.WithContext(reqCtx)
	if resolvedTenantID != "" {
		c.Set("tenant_id", resolvedTenantID)
	}
	if resolvedWorkspaceID != "" {
		c.Set("workspace_id", resolvedWorkspaceID)
	}
}
