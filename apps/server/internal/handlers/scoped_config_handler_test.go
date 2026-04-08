package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"servify/apps/server/internal/config"
	"servify/apps/server/internal/middleware"
	"servify/apps/server/internal/models"
	auditplatform "servify/apps/server/internal/platform/audit"
	platformauth "servify/apps/server/internal/platform/auth"
	"servify/apps/server/internal/platform/configscope"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newScopedConfigTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+strings.ReplaceAll(t.Name(), "/", "_")+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&models.TenantConfig{}, &models.WorkspaceConfig{}, &models.AuditLog{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func newScopedConfigTestRouter(db *gorm.DB, tenantID, workspaceID string, withAudit bool) *gin.Engine {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("principal_kind", "admin")
		c.Set("tenant_id", tenantID)
		if workspaceID != "" {
			c.Set("workspace_id", workspaceID)
		}
		c.Set("user_id", uint(9))
		c.Request = c.Request.WithContext(platformauth.ContextWithScope(context.Background(), tenantID, workspaceID))
		c.Next()
	})
	if withAudit {
		r.Use(middleware.AuditMiddleware(db))
	}
	RegisterScopedConfigRoutes(&r.RouterGroup, NewScopedConfigHandler(configscope.NewGormConfigStore(db), auditplatform.NewGormQueryService(db)))
	return r
}

func TestScopedConfigHandlerTenantRoundTrip(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := newScopedConfigTestDB(t)
	r := newScopedConfigTestRouter(db, "tenant-a", "", false)

	putBody := `{"portal":{"brand_name":"Tenant Brand"},"openai":{"base_url":"https://tenant.example/v1"},"session_risk":{"high_risk_score":9}}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/security/config/tenant", strings.NewReader(putBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("put expected 200 got %d body=%s", w.Code, w.Body.String())
	}

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/security/config/tenant", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("get expected 200 got %d body=%s", w.Code, w.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	configObj := body["config"].(map[string]any)
	portal := configObj["portal"].(map[string]any)
	if portal["brand_name"].(string) != "Tenant Brand" {
		t.Fatalf("unexpected portal payload: %+v", portal)
	}
	openai := configObj["openai"].(map[string]any)
	if openai["base_url"].(string) != "https://tenant.example/v1" {
		t.Fatalf("unexpected openai payload: %+v", openai)
	}
	sessionRisk := configObj["session_risk"].(map[string]any)
	if int(sessionRisk["high_risk_score"].(float64)) != 9 {
		t.Fatalf("unexpected session_risk payload: %+v", sessionRisk)
	}
}

func TestScopedConfigHandlerWorkspaceRoundTrip(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := newScopedConfigTestDB(t)
	r := newScopedConfigTestRouter(db, "tenant-a", "workspace-1", false)

	putBody := `{"weknora":{"enabled":true,"knowledge_base_id":"kb-workspace"}}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/security/config/workspace", strings.NewReader(putBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("put expected 200 got %d body=%s", w.Code, w.Body.String())
	}

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/security/config/workspace", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("get expected 200 got %d body=%s", w.Code, w.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	configObj := body["config"].(map[string]any)
	weknora := configObj["weknora"].(map[string]any)
	if weknora["knowledge_base_id"].(string) != "kb-workspace" {
		t.Fatalf("unexpected weknora payload: %+v", weknora)
	}
}

func TestBuildScopedConfigDiffIncludesSessionRisk(t *testing.T) {
	diff := buildScopedConfigDiff(
		&configscope.ScopedConfigDocument{
			TenantID:    "tenant-a",
			WorkspaceID: "workspace-1",
			SessionRisk: &config.SessionRiskPolicyConfig{
				HighRiskScore: 4,
			},
		},
		&configscope.ScopedConfigDocument{
			TenantID:    "tenant-a",
			WorkspaceID: "workspace-1",
			SessionRisk: &config.SessionRiskPolicyConfig{
				HighRiskScore: 8,
			},
		},
	)

	if diff["session_risk_changed"] != true {
		t.Fatalf("expected session_risk_changed=true diff=%+v", diff)
	}
	paths := diff["session_risk_paths"].([]string)
	if len(paths) != 1 || paths[0] != "session_risk.high_risk_score" {
		t.Fatalf("unexpected session_risk_paths=%v", paths)
	}
	changes := diff["session_risk_changes"].([]gin.H)
	if len(changes) != 1 {
		t.Fatalf("expected one session_risk change diff=%+v", diff)
	}
	if changes[0]["path"].(string) != "session_risk.high_risk_score" {
		t.Fatalf("unexpected session_risk change=%+v", changes[0])
	}
}

func TestScopedConfigHandlerRejectsMissingScope(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	RegisterScopedConfigRoutes(&r.RouterGroup, NewScopedConfigHandler(nil))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/security/config/tenant", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d body=%s", w.Code, w.Body.String())
	}
}

func TestScopedConfigHandlerWritesAuditLog(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := newScopedConfigTestDB(t)
	r := newScopedConfigTestRouter(db, "tenant-a", "workspace-1", true)
	h := NewScopedConfigHandler(configscope.NewGormConfigStore(db), auditplatform.NewGormQueryService(db))

	seed := configscope.ScopedConfigDocument{TenantID: "tenant-a", WorkspaceID: "workspace-1", WeKnora: &configscopeTestWeKnora}
	if _, err := h.store.UpsertWorkspaceConfig(context.Background(), "tenant-a", "workspace-1", seed); err != nil {
		t.Fatalf("seed workspace config: %v", err)
	}

	putBody := `{"weknora":{"enabled":true,"knowledge_base_id":"kb-updated","api_key":"secret-key"}}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/security/config/workspace", strings.NewReader(putBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("put expected 200 got %d body=%s", w.Code, w.Body.String())
	}

	var logs []models.AuditLog
	if err := db.Order("id asc").Find(&logs).Error; err != nil {
		t.Fatalf("load audit logs: %v", err)
	}
	if len(logs) != 1 {
		t.Fatalf("expected 1 audit log got %d", len(logs))
	}
	entry := logs[0]
	if entry.Action != "scoped_config.workspace.update" {
		t.Fatalf("action = %q want scoped_config.workspace.update", entry.Action)
	}
	if entry.ResourceType != "scoped_config" {
		t.Fatalf("resource_type = %q want scoped_config", entry.ResourceType)
	}
	if entry.ResourceID != "tenant-a/workspace-1" {
		t.Fatalf("resource_id = %q want tenant-a/workspace-1", entry.ResourceID)
	}
	if entry.ActorUserID == nil || *entry.ActorUserID != 9 {
		t.Fatalf("actor_user_id = %v want 9", entry.ActorUserID)
	}
	if entry.TenantID != "tenant-a" || entry.WorkspaceID != "workspace-1" {
		t.Fatalf("scope = %s/%s", entry.TenantID, entry.WorkspaceID)
	}
	if !strings.Contains(entry.BeforeJSON, "kb-old") {
		t.Fatalf("before_json = %q want kb-old", entry.BeforeJSON)
	}
	if !strings.Contains(entry.AfterJSON, "kb-updated") {
		t.Fatalf("after_json = %q want kb-updated", entry.AfterJSON)
	}
	if strings.Contains(entry.RequestJSON, "secret-key") || strings.Contains(entry.AfterJSON, "secret-key") {
		t.Fatalf("expected secret redaction, got request=%q after=%q", entry.RequestJSON, entry.AfterJSON)
	}
	if !strings.Contains(entry.RequestJSON, "[REDACTED]") || !strings.Contains(entry.AfterJSON, "[REDACTED]") {
		t.Fatalf("expected redaction markers, got request=%q after=%q", entry.RequestJSON, entry.AfterJSON)
	}
}

func TestScopedConfigHandlerHistoryListsScopedAuditEntries(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := newScopedConfigTestDB(t)
	r := newScopedConfigTestRouter(db, "tenant-a", "workspace-1", true)

	for _, body := range []string{
		`{"weknora":{"enabled":true,"knowledge_base_id":"kb-v1"}}`,
		`{"weknora":{"enabled":true,"knowledge_base_id":"kb-v2"}}`,
	} {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/security/config/workspace", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("put expected 200 got %d body=%s", w.Code, w.Body.String())
		}
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/security/config/workspace/history?page_size=10", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("history expected 200 got %d body=%s", w.Code, w.Body.String())
	}
	var body PaginatedResponse
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal history: %v", err)
	}
	if body.Total != 2 {
		t.Fatalf("history total = %d want 2 body=%s", body.Total, w.Body.String())
	}
	items := body.Data.([]any)
	if len(items) != 2 {
		t.Fatalf("history items = %d want 2 body=%s", len(items), w.Body.String())
	}
	first := items[0].(map[string]any)
	if first["operation"].(string) != "update" {
		t.Fatalf("operation = %v body=%s", first["operation"], w.Body.String())
	}
	if first["has_snapshot"] != true || first["can_rollback"] != true {
		t.Fatalf("unexpected rollback metadata body=%s", w.Body.String())
	}
	if first["preview_path"].(string) == "" || first["rollback_path"].(string) == "" {
		t.Fatalf("expected preview/rollback paths body=%s", w.Body.String())
	}
	audit := first["audit"].(map[string]any)
	if audit["resource_type"].(string) != "scoped_config" {
		t.Fatalf("resource_type = %v", audit["resource_type"])
	}
}

func TestScopedConfigHandlerRollbackRestoresSnapshot(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := newScopedConfigTestDB(t)
	r := newScopedConfigTestRouter(db, "tenant-a", "workspace-1", true)
	store := configscope.NewGormConfigStore(db)

	current := configscope.ScopedConfigDocument{
		TenantID:    "tenant-a",
		WorkspaceID: "workspace-1",
		WeKnora: &config.WeKnoraConfig{
			Enabled:         true,
			KnowledgeBaseID: "kb-current",
		},
	}
	if _, err := store.UpsertWorkspaceConfig(context.Background(), "tenant-a", "workspace-1", current); err != nil {
		t.Fatalf("seed current config: %v", err)
	}

	rollbackSnapshot := configscope.ScopedConfigDocument{
		TenantID:    "tenant-a",
		WorkspaceID: "workspace-1",
		WeKnora: &config.WeKnoraConfig{
			Enabled:         true,
			KnowledgeBaseID: "kb-rollback",
			APIKey:          "rollback-secret",
		},
	}
	afterJSON, err := json.Marshal(rollbackSnapshot)
	if err != nil {
		t.Fatalf("marshal snapshot: %v", err)
	}
	seedLog := models.AuditLog{
		PrincipalKind: "admin",
		Action:        "scoped_config.workspace.update",
		ResourceType:  "scoped_config",
		ResourceID:    "tenant-a/workspace-1",
		Route:         "/api/security/config/workspace",
		Method:        http.MethodPut,
		StatusCode:    http.StatusOK,
		Success:       true,
		TenantID:      "tenant-a",
		WorkspaceID:   "workspace-1",
		AfterJSON:     string(afterJSON),
	}
	if err := db.Create(&seedLog).Error; err != nil {
		t.Fatalf("seed audit log: %v", err)
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/security/config/workspace/rollback/1?confirm=true", strings.NewReader(`{"confirm":true}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("rollback expected 200 got %d body=%s", w.Code, w.Body.String())
	}

	doc, ok, err := store.GetWorkspaceConfig(context.Background(), "tenant-a", "workspace-1")
	if err != nil || !ok {
		t.Fatalf("load rolled back config ok=%v err=%v", ok, err)
	}
	if doc.WeKnora == nil || doc.WeKnora.KnowledgeBaseID != "kb-rollback" {
		t.Fatalf("rolled back config = %+v", doc)
	}

	var logs []models.AuditLog
	if err := db.Order("id asc").Find(&logs).Error; err != nil {
		t.Fatalf("load audit logs: %v", err)
	}
	if len(logs) != 2 {
		t.Fatalf("expected 2 audit logs got %d", len(logs))
	}
	rollbackLog := logs[1]
	if rollbackLog.Action != "scoped_config.workspace.rollback" {
		t.Fatalf("rollback action = %q want scoped_config.workspace.rollback", rollbackLog.Action)
	}
	if !strings.Contains(rollbackLog.BeforeJSON, "kb-current") {
		t.Fatalf("rollback before_json = %q", rollbackLog.BeforeJSON)
	}
	if !strings.Contains(rollbackLog.AfterJSON, "kb-rollback") {
		t.Fatalf("rollback after_json = %q", rollbackLog.AfterJSON)
	}
	if strings.Contains(rollbackLog.AfterJSON, "rollback-secret") {
		t.Fatalf("expected redacted rollback secret, got %q", rollbackLog.AfterJSON)
	}
}

var configscopeTestWeKnora = config.WeKnoraConfig{
	Enabled:         true,
	KnowledgeBaseID: "kb-old",
	APIKey:          "seed-secret",
}

func TestScopedConfigHandlerHistoryEntryShowsDiffPreview(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := newScopedConfigTestDB(t)
	r := newScopedConfigTestRouter(db, "tenant-a", "workspace-1", false)
	store := configscope.NewGormConfigStore(db)

	current := configscope.ScopedConfigDocument{
		TenantID:    "tenant-a",
		WorkspaceID: "workspace-1",
		OpenAI:      &config.OpenAIConfig{Model: "gpt-current"},
		WeKnora:     &config.WeKnoraConfig{Enabled: true, KnowledgeBaseID: "kb-current"},
	}
	if _, err := store.UpsertWorkspaceConfig(context.Background(), "tenant-a", "workspace-1", current); err != nil {
		t.Fatalf("seed current config: %v", err)
	}

	snapshot := configscope.ScopedConfigDocument{
		TenantID:    "tenant-a",
		WorkspaceID: "workspace-1",
		OpenAI:      &config.OpenAIConfig{Model: "gpt-preview"},
		WeKnora:     &config.WeKnoraConfig{Enabled: true, KnowledgeBaseID: "kb-current"},
	}
	afterJSON, err := json.Marshal(snapshot)
	if err != nil {
		t.Fatalf("marshal snapshot: %v", err)
	}
	seedLog := models.AuditLog{
		PrincipalKind: "admin",
		Action:        "scoped_config.workspace.update",
		ResourceType:  "scoped_config",
		ResourceID:    "tenant-a/workspace-1",
		Route:         "/api/security/config/workspace",
		Method:        http.MethodPut,
		StatusCode:    http.StatusOK,
		Success:       true,
		TenantID:      "tenant-a",
		WorkspaceID:   "workspace-1",
		AfterJSON:     string(afterJSON),
	}
	if err := db.Create(&seedLog).Error; err != nil {
		t.Fatalf("seed audit log: %v", err)
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/security/config/workspace/history/1", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("history entry expected 200 got %d body=%s", w.Code, w.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}
	diff := body["diff"].(map[string]any)
	if diff["openai_changed"] != true {
		t.Fatalf("expected openai_changed=true body=%s", w.Body.String())
	}
	if diff["weknora_changed"] != false {
		t.Fatalf("expected weknora_changed=false body=%s", w.Body.String())
	}
	changedPaths := diff["changed_paths"].([]any)
	if len(changedPaths) != 1 || changedPaths[0].(string) != "openai.model" {
		t.Fatalf("expected changed_paths=[openai.model] body=%s", w.Body.String())
	}
	openaiPaths := diff["openai_paths"].([]any)
	if len(openaiPaths) != 1 || openaiPaths[0].(string) != "openai.model" {
		t.Fatalf("expected openai_paths=[openai.model] body=%s", w.Body.String())
	}
	changes := diff["changes"].([]any)
	if len(changes) != 1 {
		t.Fatalf("expected one structured change body=%s", w.Body.String())
	}
	change := changes[0].(map[string]any)
	if change["path"].(string) != "openai.model" {
		t.Fatalf("expected change path=openai.model body=%s", w.Body.String())
	}
	if change["type"].(string) != "updated" {
		t.Fatalf("expected change type=updated body=%s", w.Body.String())
	}
	if change["current"].(string) != "gpt-current" || change["snapshot"].(string) != "gpt-preview" {
		t.Fatalf("unexpected change payload body=%s", w.Body.String())
	}
	openaiChanges := diff["openai_changes"].([]any)
	if len(openaiChanges) != 1 {
		t.Fatalf("expected one openai change body=%s", w.Body.String())
	}
	snapshotBody := body["snapshot"].(map[string]any)
	openai := snapshotBody["openai"].(map[string]any)
	if openai["model"].(string) != "gpt-preview" {
		t.Fatalf("snapshot openai model=%v", openai["model"])
	}
}

func TestScopedConfigHandlerHistoryEntryMarksAddedFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := newScopedConfigTestDB(t)
	r := newScopedConfigTestRouter(db, "tenant-a", "workspace-1", false)

	snapshot := configscope.ScopedConfigDocument{
		TenantID:    "tenant-a",
		WorkspaceID: "workspace-1",
		Portal:      &config.PortalConfig{BrandName: "Added Brand"},
	}
	afterJSON, err := json.Marshal(snapshot)
	if err != nil {
		t.Fatalf("marshal snapshot: %v", err)
	}
	seedLog := models.AuditLog{
		PrincipalKind: "admin",
		Action:        "scoped_config.workspace.update",
		ResourceType:  "scoped_config",
		ResourceID:    "tenant-a/workspace-1",
		Route:         "/api/security/config/workspace",
		Method:        http.MethodPut,
		StatusCode:    http.StatusOK,
		Success:       true,
		TenantID:      "tenant-a",
		WorkspaceID:   "workspace-1",
		AfterJSON:     string(afterJSON),
	}
	if err := db.Create(&seedLog).Error; err != nil {
		t.Fatalf("seed audit log: %v", err)
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/security/config/workspace/history/1", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("history entry expected 200 got %d body=%s", w.Code, w.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}
	diff := body["diff"].(map[string]any)
	changes := diff["changes"].([]any)
	if len(changes) != 1 {
		t.Fatalf("expected one added change body=%s", w.Body.String())
	}
	change := changes[0].(map[string]any)
	if change["path"].(string) != "portal.brand_name" {
		t.Fatalf("unexpected path body=%s", w.Body.String())
	}
	if change["type"].(string) != "added" {
		t.Fatalf("expected added change body=%s", w.Body.String())
	}
	if change["current"] != nil || change["snapshot"].(string) != "Added Brand" {
		t.Fatalf("unexpected added payload body=%s", w.Body.String())
	}
}

func TestScopedConfigHandlerRollbackRequiresConfirmation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := newScopedConfigTestDB(t)
	r := newScopedConfigTestRouter(db, "tenant-a", "workspace-1", true)

	snapshot := configscope.ScopedConfigDocument{
		TenantID:    "tenant-a",
		WorkspaceID: "workspace-1",
		WeKnora:     &config.WeKnoraConfig{Enabled: true, KnowledgeBaseID: "kb-rollback"},
	}
	afterJSON, err := json.Marshal(snapshot)
	if err != nil {
		t.Fatalf("marshal snapshot: %v", err)
	}
	seedLog := models.AuditLog{
		PrincipalKind: "admin",
		Action:        "scoped_config.workspace.update",
		ResourceType:  "scoped_config",
		ResourceID:    "tenant-a/workspace-1",
		Route:         "/api/security/config/workspace",
		Method:        http.MethodPut,
		StatusCode:    http.StatusOK,
		Success:       true,
		TenantID:      "tenant-a",
		WorkspaceID:   "workspace-1",
		AfterJSON:     string(afterJSON),
	}
	if err := db.Create(&seedLog).Error; err != nil {
		t.Fatalf("seed audit log: %v", err)
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/security/config/workspace/rollback/1", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d body=%s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "confirm=true") {
		t.Fatalf("expected confirm hint body=%s", w.Body.String())
	}
}
