package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"servify/apps/server/internal/models"
	auditplatform "servify/apps/server/internal/platform/audit"
	platformauth "servify/apps/server/internal/platform/auth"
	"servify/apps/server/internal/platform/configscope"

	"github.com/gin-gonic/gin"
)

type ScopedConfigHandler struct {
	store *configscope.GormConfigStore
	audit auditplatform.QueryService
}

func NewScopedConfigHandler(store *configscope.GormConfigStore, auditServices ...auditplatform.QueryService) *ScopedConfigHandler {
	h := &ScopedConfigHandler{store: store}
	if len(auditServices) > 0 {
		h.audit = auditServices[0]
	}
	return h
}

func (h *ScopedConfigHandler) GetTenantConfig(c *gin.Context) {
	tenantID := platformauth.TenantIDFromContext(c.Request.Context())
	if tenantID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Tenant scope required", Message: "tenant_id is required in request scope"})
		return
	}
	if h == nil || h.store == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Config store unavailable", Message: "config store not configured"})
		return
	}
	doc, ok, err := h.store.GetTenantConfig(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to load tenant config", Message: err.Error()})
		return
	}
	if !ok {
		c.JSON(http.StatusOK, gin.H{"tenant_id": tenantID, "config": configscope.ScopedConfigDocument{TenantID: tenantID}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tenant_id": tenantID, "config": doc})
}

func (h *ScopedConfigHandler) PutTenantConfig(c *gin.Context) {
	tenantID := platformauth.TenantIDFromContext(c.Request.Context())
	if tenantID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Tenant scope required", Message: "tenant_id is required in request scope"})
		return
	}
	if h == nil || h.store == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Config store unavailable", Message: "config store not configured"})
		return
	}
	var payload configscope.ScopedConfigDocument
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request", Message: err.Error()})
		return
	}
	if before, ok, err := h.store.GetTenantConfig(c.Request.Context(), tenantID); err == nil && ok {
		auditplatform.SetBefore(c, before)
	}
	setScopedConfigAuditMeta(c, "tenant", tenantID, "", "update")
	doc, err := h.store.UpsertTenantConfig(c.Request.Context(), tenantID, payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Failed to save tenant config", Message: err.Error()})
		return
	}
	auditplatform.SetAfter(c, doc)
	c.JSON(http.StatusOK, gin.H{"tenant_id": tenantID, "config": doc})
}

func (h *ScopedConfigHandler) GetTenantConfigHistory(c *gin.Context) {
	h.listHistory(c, "tenant")
}

func (h *ScopedConfigHandler) GetTenantConfigHistoryEntry(c *gin.Context) {
	h.historyEntry(c, "tenant")
}

func (h *ScopedConfigHandler) RollbackTenantConfig(c *gin.Context) {
	h.rollback(c, "tenant")
}

func (h *ScopedConfigHandler) GetWorkspaceConfig(c *gin.Context) {
	tenantID := platformauth.TenantIDFromContext(c.Request.Context())
	workspaceID := platformauth.WorkspaceIDFromContext(c.Request.Context())
	if tenantID == "" || workspaceID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Workspace scope required", Message: "tenant_id and workspace_id are required in request scope"})
		return
	}
	if h == nil || h.store == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Config store unavailable", Message: "config store not configured"})
		return
	}
	doc, ok, err := h.store.GetWorkspaceConfig(c.Request.Context(), tenantID, workspaceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to load workspace config", Message: err.Error()})
		return
	}
	if !ok {
		c.JSON(http.StatusOK, gin.H{"tenant_id": tenantID, "workspace_id": workspaceID, "config": configscope.ScopedConfigDocument{TenantID: tenantID, WorkspaceID: workspaceID}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tenant_id": tenantID, "workspace_id": workspaceID, "config": doc})
}

func (h *ScopedConfigHandler) PutWorkspaceConfig(c *gin.Context) {
	tenantID := platformauth.TenantIDFromContext(c.Request.Context())
	workspaceID := platformauth.WorkspaceIDFromContext(c.Request.Context())
	if tenantID == "" || workspaceID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Workspace scope required", Message: "tenant_id and workspace_id are required in request scope"})
		return
	}
	if h == nil || h.store == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Config store unavailable", Message: "config store not configured"})
		return
	}
	var payload configscope.ScopedConfigDocument
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request", Message: err.Error()})
		return
	}
	if before, ok, err := h.store.GetWorkspaceConfig(c.Request.Context(), tenantID, workspaceID); err == nil && ok {
		auditplatform.SetBefore(c, before)
	}
	setScopedConfigAuditMeta(c, "workspace", tenantID, workspaceID, "update")
	doc, err := h.store.UpsertWorkspaceConfig(c.Request.Context(), tenantID, workspaceID, payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Failed to save workspace config", Message: err.Error()})
		return
	}
	auditplatform.SetAfter(c, doc)
	c.JSON(http.StatusOK, gin.H{"tenant_id": tenantID, "workspace_id": workspaceID, "config": doc})
}

func (h *ScopedConfigHandler) GetWorkspaceConfigHistory(c *gin.Context) {
	h.listHistory(c, "workspace")
}

func (h *ScopedConfigHandler) GetWorkspaceConfigHistoryEntry(c *gin.Context) {
	h.historyEntry(c, "workspace")
}

func (h *ScopedConfigHandler) RollbackWorkspaceConfig(c *gin.Context) {
	h.rollback(c, "workspace")
}

func (h *ScopedConfigHandler) historyEntry(c *gin.Context, scope string) {
	tenantID, workspaceID, ok := requireScope(c, scope)
	if !ok {
		return
	}
	if h == nil || h.store == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Config store unavailable", Message: "config store not configured"})
		return
	}
	entry, snapshot, response, err := h.loadHistoryEntry(c, scope, tenantID, workspaceID)
	if response != nil {
		c.JSON(response.Code, response)
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to load config history entry", Message: err.Error()})
		return
	}
	current, err := h.currentConfigForScope(c.Request.Context(), scope, tenantID, workspaceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to load current scoped config", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"audit":    entry,
		"snapshot": snapshot,
		"current":  current,
		"diff":     buildScopedConfigDiff(current, snapshot),
	})
}

func (h *ScopedConfigHandler) listHistory(c *gin.Context, scope string) {
	tenantID, workspaceID, ok := requireScope(c, scope)
	if !ok {
		return
	}
	if h == nil || h.audit == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Audit service unavailable", Message: "audit query service not configured"})
		return
	}

	query := auditplatform.ListQuery{
		Action:       c.Query("action"),
		ResourceType: "scoped_config",
		ResourceID:   scopedConfigResourceID(tenantID, workspaceID),
		TenantID:     tenantID,
		WorkspaceID:  workspaceID,
		Page:         intQuery(c, "page", 1),
		PageSize:     intQuery(c, "page_size", 20),
	}
	items, total, err := h.audit.List(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to list config history", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, PaginatedResponse{Data: buildScopedConfigHistoryItems(scope, items), Total: total, Page: query.Page, PageSize: query.PageSize})
}

func buildScopedConfigHistoryItems(scope string, items []models.AuditLog) []gin.H {
	result := make([]gin.H, 0, len(items))
	for _, item := range items {
		hasSnapshot := strings.TrimSpace(item.AfterJSON) != ""
		operation := scopedConfigOperation(item.Action)
		result = append(result, gin.H{
			"audit":         item,
			"operation":     operation,
			"has_snapshot":  hasSnapshot,
			"can_rollback":  hasSnapshot && (operation == "update" || operation == "rollback"),
			"preview_path":  fmt.Sprintf("/security/config/%s/history/%d", scope, item.ID),
			"rollback_path": fmt.Sprintf("/security/config/%s/rollback/%d", scope, item.ID),
		})
	}
	return result
}

func scopedConfigOperation(action string) string {
	action = strings.TrimSpace(action)
	parts := strings.Split(action, ".")
	if len(parts) >= 3 {
		return parts[len(parts)-1]
	}
	return action
}

type scopedConfigRollbackRequest struct {
	Confirm bool `json:"confirm"`
}

func requireRollbackConfirmation(c *gin.Context) bool {
	if c == nil {
		return false
	}
	if strings.EqualFold(strings.TrimSpace(c.Query("confirm")), "true") {
		return true
	}
	if c.Request == nil || c.Request.Body == nil {
		return false
	}
	var payload scopedConfigRollbackRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		return false
	}
	return payload.Confirm
}

func (h *ScopedConfigHandler) rollback(c *gin.Context, scope string) {
	tenantID, workspaceID, ok := requireScope(c, scope)
	if !ok {
		return
	}
	if h == nil || h.store == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Config store unavailable", Message: "config store not configured"})
		return
	}
	if h.audit == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Audit service unavailable", Message: "audit query service not configured"})
		return
	}
	if !requireRollbackConfirmation(c) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Rollback confirmation required", Message: "set confirm=true in query or JSON body to execute rollback"})
		return
	}

	auditID, err := strconv.ParseUint(strings.TrimSpace(c.Param("audit_id")), 10, 32)
	if err != nil || auditID == 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid audit_id", Message: "audit_id must be a positive integer"})
		return
	}
	entry, snapshot, response, err := h.loadHistoryEntry(c, scope, tenantID, workspaceID)
	if response != nil {
		c.JSON(response.Code, response)
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to load config history entry", Message: err.Error()})
		return
	}
	if scope == "tenant" {
		snapshot.TenantID = tenantID
		snapshot.WorkspaceID = ""
		if before, ok, err := h.store.GetTenantConfig(c.Request.Context(), tenantID); err == nil && ok {
			auditplatform.SetBefore(c, before)
		}
		setScopedConfigAuditMeta(c, scope, tenantID, "", "rollback")
		doc, err := h.store.UpsertTenantConfig(c.Request.Context(), tenantID, *snapshot)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Failed to rollback tenant config", Message: err.Error()})
			return
		}
		auditplatform.SetAfter(c, doc)
		c.JSON(http.StatusOK, gin.H{"tenant_id": tenantID, "source_audit_id": entry.ID, "config": doc})
		return
	}

	snapshot.TenantID = tenantID
	snapshot.WorkspaceID = workspaceID
	if before, ok, err := h.store.GetWorkspaceConfig(c.Request.Context(), tenantID, workspaceID); err == nil && ok {
		auditplatform.SetBefore(c, before)
	}
	setScopedConfigAuditMeta(c, scope, tenantID, workspaceID, "rollback")
	doc, err := h.store.UpsertWorkspaceConfig(c.Request.Context(), tenantID, workspaceID, *snapshot)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Failed to rollback workspace config", Message: err.Error()})
		return
	}
	auditplatform.SetAfter(c, doc)
	c.JSON(http.StatusOK, gin.H{"tenant_id": tenantID, "workspace_id": workspaceID, "source_audit_id": entry.ID, "config": doc})
}

func requireScope(c *gin.Context, scope string) (string, string, bool) {
	tenantID := platformauth.TenantIDFromContext(c.Request.Context())
	workspaceID := platformauth.WorkspaceIDFromContext(c.Request.Context())
	if scope == "tenant" {
		if tenantID == "" {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Tenant scope required", Message: "tenant_id is required in request scope"})
			return "", "", false
		}
		return tenantID, "", true
	}
	if tenantID == "" || workspaceID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Workspace scope required", Message: "tenant_id and workspace_id are required in request scope"})
		return "", "", false
	}
	return tenantID, workspaceID, true
}

func setScopedConfigAuditMeta(c *gin.Context, scope, tenantID, workspaceID, operation string) {
	if c == nil {
		return
	}
	auditplatform.SetAction(c, fmt.Sprintf("scoped_config.%s.%s", scope, operation))
	auditplatform.SetResourceType(c, "scoped_config")
	auditplatform.SetResourceID(c, scopedConfigResourceID(tenantID, workspaceID))
}

func scopedConfigResourceID(tenantID, workspaceID string) string {
	if workspaceID == "" {
		return tenantID
	}
	return tenantID + "/" + workspaceID
}

func RegisterScopedConfigRoutes(r *gin.RouterGroup, handler *ScopedConfigHandler) {
	if r == nil || handler == nil {
		return
	}
	security := r.Group("/security/config")
	{
		security.GET("/tenant", handler.GetTenantConfig)
		security.PUT("/tenant", handler.PutTenantConfig)
		security.GET("/tenant/history", handler.GetTenantConfigHistory)
		security.GET("/tenant/history/:audit_id", handler.GetTenantConfigHistoryEntry)
		security.POST("/tenant/rollback/:audit_id", handler.RollbackTenantConfig)
		security.GET("/workspace", handler.GetWorkspaceConfig)
		security.PUT("/workspace", handler.PutWorkspaceConfig)
		security.GET("/workspace/history", handler.GetWorkspaceConfigHistory)
		security.GET("/workspace/history/:audit_id", handler.GetWorkspaceConfigHistoryEntry)
		security.POST("/workspace/rollback/:audit_id", handler.RollbackWorkspaceConfig)
	}
}

func (h *ScopedConfigHandler) loadHistoryEntry(c *gin.Context, scope, tenantID, workspaceID string) (*models.AuditLog, *configscope.ScopedConfigDocument, *ErrorResponse, error) {
	if h == nil || h.audit == nil {
		return nil, nil, &ErrorResponse{Code: http.StatusNotFound, Error: "Audit service unavailable", Message: "audit query service not configured"}, nil
	}
	auditID, err := strconv.ParseUint(strings.TrimSpace(c.Param("audit_id")), 10, 32)
	if err != nil || auditID == 0 {
		return nil, nil, &ErrorResponse{Code: http.StatusBadRequest, Error: "Invalid audit_id", Message: "audit_id must be a positive integer"}, nil
	}
	resourceID := scopedConfigResourceID(tenantID, workspaceID)
	entry, err := h.audit.Get(c.Request.Context(), uint(auditID), auditplatform.QueryScope{TenantID: tenantID, WorkspaceID: workspaceID})
	if err != nil {
		return nil, nil, nil, err
	}
	if entry == nil || entry.ResourceType != "scoped_config" || entry.ResourceID != resourceID {
		return nil, nil, &ErrorResponse{Code: http.StatusNotFound, Error: "Config history entry not found", Message: "no matching scoped config audit entry found"}, nil
	}
	if strings.TrimSpace(entry.AfterJSON) == "" {
		return nil, nil, &ErrorResponse{Code: http.StatusBadRequest, Error: "Rollback snapshot unavailable", Message: "selected audit entry does not contain an after snapshot"}, nil
	}
	var snapshot configscope.ScopedConfigDocument
	if err := json.Unmarshal([]byte(entry.AfterJSON), &snapshot); err != nil {
		return nil, nil, &ErrorResponse{Code: http.StatusBadRequest, Error: "Invalid rollback snapshot", Message: err.Error()}, nil
	}
	if snapshot.TenantID != "" && snapshot.TenantID != tenantID {
		return nil, nil, &ErrorResponse{Code: http.StatusBadRequest, Error: "Rollback scope mismatch", Message: "snapshot tenant scope does not match request scope"}, nil
	}
	if scope == "workspace" && snapshot.WorkspaceID != "" && snapshot.WorkspaceID != workspaceID {
		return nil, nil, &ErrorResponse{Code: http.StatusBadRequest, Error: "Rollback scope mismatch", Message: "snapshot workspace scope does not match request scope"}, nil
	}
	return entry, &snapshot, nil, nil
}

func (h *ScopedConfigHandler) currentConfigForScope(ctx context.Context, scope, tenantID, workspaceID string) (*configscope.ScopedConfigDocument, error) {
	if h == nil || h.store == nil {
		return nil, nil
	}
	if scope == "tenant" {
		if doc, ok, err := h.store.GetTenantConfig(ctx, tenantID); err != nil {
			return nil, err
		} else if ok {
			return doc, nil
		}
		return &configscope.ScopedConfigDocument{TenantID: tenantID}, nil
	}
	if doc, ok, err := h.store.GetWorkspaceConfig(ctx, tenantID, workspaceID); err != nil {
		return nil, err
	} else if ok {
		return doc, nil
	}
	return &configscope.ScopedConfigDocument{TenantID: tenantID, WorkspaceID: workspaceID}, nil
}

func buildScopedConfigDiff(current, snapshot *configscope.ScopedConfigDocument) gin.H {
	if current == nil {
		current = &configscope.ScopedConfigDocument{}
	}
	if snapshot == nil {
		snapshot = &configscope.ScopedConfigDocument{}
	}
	portalChanges := diffJSONChanges(current.Portal, snapshot.Portal, "portal")
	openAIChanges := diffJSONChanges(current.OpenAI, snapshot.OpenAI, "openai")
	weKnoraChanges := diffJSONChanges(current.WeKnora, snapshot.WeKnora, "weknora")
	changes := append(append([]gin.H{}, portalChanges...), openAIChanges...)
	changes = append(changes, weKnoraChanges...)
	changedPaths := changePaths(changes)
	return gin.H{
		"portal_changed":    len(portalChanges) > 0,
		"openai_changed":    len(openAIChanges) > 0,
		"weknora_changed":   len(weKnoraChanges) > 0,
		"scope_changed":     current.TenantID != snapshot.TenantID || current.WorkspaceID != snapshot.WorkspaceID,
		"current_sections":  presentSections(current),
		"snapshot_sections": presentSections(snapshot),
		"changed_paths":     changedPaths,
		"portal_paths":      changePaths(portalChanges),
		"openai_paths":      changePaths(openAIChanges),
		"weknora_paths":     changePaths(weKnoraChanges),
		"changes":           changes,
		"portal_changes":    portalChanges,
		"openai_changes":    openAIChanges,
		"weknora_changes":   weKnoraChanges,
	}
}

func presentSections(doc *configscope.ScopedConfigDocument) []string {
	if doc == nil {
		return nil
	}
	sections := make([]string, 0, 3)
	if doc.Portal != nil {
		sections = append(sections, "portal")
	}
	if doc.OpenAI != nil {
		sections = append(sections, "openai")
	}
	if doc.WeKnora != nil {
		sections = append(sections, "weknora")
	}
	return sections
}

func jsonValue(value interface{}) string {
	if value == nil {
		return ""
	}
	body, err := json.Marshal(value)
	if err != nil {
		return ""
	}
	return string(body)
}

func diffJSONChanges(current, snapshot interface{}, prefix string) []gin.H {
	var currentValue interface{}
	var snapshotValue interface{}
	if body, err := json.Marshal(current); err == nil && len(body) > 0 && string(body) != "null" {
		_ = json.Unmarshal(body, &currentValue)
	}
	if body, err := json.Marshal(snapshot); err == nil && len(body) > 0 && string(body) != "null" {
		_ = json.Unmarshal(body, &snapshotValue)
	}
	changes := collectDiffChanges(currentValue, snapshotValue, prefix)
	sort.Slice(changes, func(i, j int) bool {
		return changes[i]["path"].(string) < changes[j]["path"].(string)
	})
	return changes
}

func collectDiffChanges(current, snapshot interface{}, prefix string) []gin.H {
	if valuesEqual(current, snapshot) {
		return nil
	}
	currentMap, currentIsMap := current.(map[string]interface{})
	snapshotMap, snapshotIsMap := snapshot.(map[string]interface{})
	if currentIsMap || snapshotIsMap {
		keys := make(map[string]struct{}, len(currentMap)+len(snapshotMap))
		for key := range currentMap {
			keys[key] = struct{}{}
		}
		for key := range snapshotMap {
			keys[key] = struct{}{}
		}
		ordered := make([]string, 0, len(keys))
		for key := range keys {
			ordered = append(ordered, key)
		}
		sort.Strings(ordered)
		var changes []gin.H
		for _, key := range ordered {
			nextPrefix := key
			if prefix != "" {
				nextPrefix = prefix + "." + key
			}
			changes = append(changes, collectDiffChanges(currentMap[key], snapshotMap[key], nextPrefix)...)
		}
		if len(changes) == 0 && prefix != "" {
			return []gin.H{{"path": prefix, "type": diffChangeType(current, snapshot), "current": current, "snapshot": snapshot}}
		}
		return changes
	}
	return []gin.H{{"path": prefix, "type": diffChangeType(current, snapshot), "current": current, "snapshot": snapshot}}
}

func diffChangeType(current, snapshot interface{}) string {
	if current == nil && snapshot != nil {
		return "added"
	}
	if current != nil && snapshot == nil {
		return "removed"
	}
	return "updated"
}

func valuesEqual(a, b interface{}) bool {
	left, _ := json.Marshal(a)
	right, _ := json.Marshal(b)
	return string(left) == string(right)
}

func changePaths(changes []gin.H) []string {
	paths := make([]string, 0, len(changes))
	for _, change := range changes {
		path, _ := change["path"].(string)
		if path != "" {
			paths = append(paths, path)
		}
	}
	sort.Strings(paths)
	return paths
}
