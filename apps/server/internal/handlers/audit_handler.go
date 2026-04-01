package handlers

import (
	"net/http"
	"strconv"
	"time"

	auditplatform "servify/apps/server/internal/platform/audit"
	platformauth "servify/apps/server/internal/platform/auth"

	"github.com/gin-gonic/gin"
)

type AuditHandler struct {
	service auditplatform.QueryService
}

func NewAuditHandler(service auditplatform.QueryService) *AuditHandler {
	return &AuditHandler{service: service}
}

func (h *AuditHandler) List(c *gin.Context) {
	if h == nil || h.service == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Audit service unavailable", Message: "audit query service not configured"})
		return
	}

	query := auditplatform.ListQuery{
		Action:        c.Query("action"),
		ResourceType:  c.Query("resource_type"),
		ResourceID:    c.Query("resource_id"),
		PrincipalKind: c.Query("principal_kind"),
		TenantID:      platformauth.TenantIDFromContext(c.Request.Context()),
		WorkspaceID:   platformauth.WorkspaceIDFromContext(c.Request.Context()),
		Page:          intQuery(c, "page", 1),
		PageSize:      intQuery(c, "page_size", 20),
	}

	if v := c.Query("actor_user_id"); v != "" {
		id, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid actor_user_id", Message: err.Error()})
			return
		}
		uid := uint(id)
		query.ActorUserID = &uid
	}

	if v := c.Query("success"); v != "" {
		switch v {
		case "true":
			b := true
			query.Success = &b
		case "false":
			b := false
			query.Success = &b
		default:
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid success", Message: "use true or false"})
			return
		}
	}

	if v := c.Query("from"); v != "" {
		tm, err := time.Parse(time.RFC3339, v)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid from", Message: err.Error()})
			return
		}
		query.From = &tm
	}
	if v := c.Query("to"); v != "" {
		tm, err := time.Parse(time.RFC3339, v)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid to", Message: err.Error()})
			return
		}
		query.To = &tm
	}

	items, total, err := h.service.List(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to list audit logs", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Data:     items,
		Total:    total,
		Page:     query.Page,
		PageSize: query.PageSize,
	})
}

func RegisterAuditRoutes(r *gin.RouterGroup, handler *AuditHandler) {
	if r == nil || handler == nil {
		return
	}
	audit := r.Group("/audit")
	{
		audit.GET("/logs", handler.List)
	}
}

func intQuery(c *gin.Context, key string, fallback int) int {
	if c == nil {
		return fallback
	}
	v := c.Query(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return fallback
	}
	return n
}
