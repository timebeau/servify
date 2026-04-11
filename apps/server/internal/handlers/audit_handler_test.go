package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"servify/apps/server/internal/models"
	auditplatform "servify/apps/server/internal/platform/audit"
	platformauth "servify/apps/server/internal/platform/auth"

	"github.com/gin-gonic/gin"
)

type stubAuditQueryService struct {
	items []models.AuditLog
	total int64
	err   error
	query auditplatform.ListQuery
	getID uint
	scope auditplatform.QueryScope
}

func (s *stubAuditQueryService) List(_ context.Context, query auditplatform.ListQuery) ([]models.AuditLog, int64, error) {
	s.query = query
	return s.items, s.total, s.err
}

func (s *stubAuditQueryService) Get(_ context.Context, id uint, scope auditplatform.QueryScope) (*models.AuditLog, error) {
	s.getID = id
	s.scope = scope
	for i := range s.items {
		if s.items[i].ID == id {
			return &s.items[i], s.err
		}
	}
	return nil, s.err
}

func TestAuditHandlerList(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &stubAuditQueryService{
		items: []models.AuditLog{{ID: 1, Action: "tickets.create"}},
		total: 1,
	}
	r := gin.New()
	RegisterAuditRoutes(&r.RouterGroup, NewAuditHandler(svc))

	req := httptest.NewRequest(http.MethodGet, "/audit/logs?action=tickets.create&principal_kind=agent&success=true&from="+time.Now().UTC().Format(time.RFC3339), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d body=%s", w.Code, w.Body.String())
	}
	if svc.query.Action != "tickets.create" || svc.query.PrincipalKind != "agent" {
		t.Fatalf("unexpected query: %+v", svc.query)
	}
	if svc.query.Success == nil || !*svc.query.Success {
		t.Fatalf("expected success filter true, got %+v", svc.query.Success)
	}
}

func TestAuditHandlerListProjectsScopeFromContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &stubAuditQueryService{
		items: []models.AuditLog{{ID: 1, Action: "tickets.create"}},
		total: 1,
	}
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Request = c.Request.WithContext(platformauth.ContextWithScope(c.Request.Context(), "tenant-a", "workspace-1"))
		c.Next()
	})
	RegisterAuditRoutes(&r.RouterGroup, NewAuditHandler(svc))

	req := httptest.NewRequest(http.MethodGet, "/audit/logs", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d body=%s", w.Code, w.Body.String())
	}
	if svc.query.TenantID != "tenant-a" || svc.query.WorkspaceID != "workspace-1" {
		t.Fatalf("unexpected scope query: %+v", svc.query)
	}
}

func TestAuditHandlerGet(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &stubAuditQueryService{
		items: []models.AuditLog{{ID: 7, Action: "tickets.create", TenantID: "tenant-a", WorkspaceID: "workspace-1"}},
	}
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Request = c.Request.WithContext(platformauth.ContextWithScope(c.Request.Context(), "tenant-a", "workspace-1"))
		c.Next()
	})
	RegisterAuditRoutes(&r.RouterGroup, NewAuditHandler(svc))

	req := httptest.NewRequest(http.MethodGet, "/audit/logs/7", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d body=%s", w.Code, w.Body.String())
	}
	if svc.getID != 7 {
		t.Fatalf("unexpected get id: %d", svc.getID)
	}
	if svc.scope.TenantID != "tenant-a" || svc.scope.WorkspaceID != "workspace-1" {
		t.Fatalf("unexpected scope: %+v", svc.scope)
	}
}

func TestAuditHandlerGetRejectsInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &stubAuditQueryService{}
	r := gin.New()
	RegisterAuditRoutes(&r.RouterGroup, NewAuditHandler(svc))

	req := httptest.NewRequest(http.MethodGet, "/audit/logs/not-a-number", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d body=%s", w.Code, w.Body.String())
	}
}

func TestAuditHandlerGetDiff(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &stubAuditQueryService{
		items: []models.AuditLog{{
			ID:          9,
			Action:      "tickets.update",
			TenantID:    "tenant-a",
			WorkspaceID: "workspace-1",
			BeforeJSON:  `{"status":"open","agent_id":null,"title":"before"}`,
			AfterJSON:   `{"status":"resolved","agent_id":42,"title":"after"}`,
		}},
	}
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Request = c.Request.WithContext(platformauth.ContextWithScope(c.Request.Context(), "tenant-a", "workspace-1"))
		c.Next()
	})
	RegisterAuditRoutes(&r.RouterGroup, NewAuditHandler(svc))

	req := httptest.NewRequest(http.MethodGet, "/audit/logs/9/diff", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d body=%s", w.Code, w.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v body=%s", err, w.Body.String())
	}
	diff := body["diff"].(map[string]any)
	if diff["changed"] != true {
		t.Fatalf("expected changed=true body=%s", w.Body.String())
	}
	if int(diff["change_count"].(float64)) < 2 {
		t.Fatalf("expected at least 2 changes body=%s", w.Body.String())
	}
	paths := diff["changed_paths"].([]any)
	if len(paths) == 0 {
		t.Fatalf("expected changed paths body=%s", w.Body.String())
	}
	if svc.scope.TenantID != "tenant-a" || svc.scope.WorkspaceID != "workspace-1" {
		t.Fatalf("unexpected scope: %+v", svc.scope)
	}
}

func TestAuditHandlerGetDiffRequiresSnapshots(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &stubAuditQueryService{
		items: []models.AuditLog{{ID: 10, Action: "tickets.update"}},
	}
	r := gin.New()
	RegisterAuditRoutes(&r.RouterGroup, NewAuditHandler(svc))

	req := httptest.NewRequest(http.MethodGet, "/audit/logs/10/diff", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d body=%s", w.Code, w.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v body=%s", err, w.Body.String())
	}
	diff := body["diff"].(map[string]any)
	if diff["changed"] != false || int(diff["change_count"].(float64)) != 0 {
		t.Fatalf("unexpected diff body=%s", w.Body.String())
	}
}

func TestAuditHandlerExportCSV(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &stubAuditQueryService{
		items: []models.AuditLog{{
			ID:            11,
			Action:        "tickets.update",
			ResourceType:  "tickets",
			ResourceID:    "1",
			PrincipalKind: "admin",
			TenantID:      "tenant-a",
			WorkspaceID:   "workspace-1",
			RequestJSON:   `{"status":"resolved"}`,
			BeforeJSON:    `{"status":"open"}`,
			AfterJSON:     `{"status":"resolved"}`,
			CreatedAt:     time.Date(2026, 4, 3, 12, 0, 0, 0, time.UTC),
		}},
		total: 1,
	}
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Request = c.Request.WithContext(platformauth.ContextWithScope(c.Request.Context(), "tenant-a", "workspace-1"))
		c.Next()
	})
	RegisterAuditRoutes(&r.RouterGroup, NewAuditHandler(svc))

	req := httptest.NewRequest(http.MethodGet, "/audit/logs/export?action=tickets.update&limit=10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d body=%s", w.Code, w.Body.String())
	}
	if contentType := w.Header().Get("Content-Type"); !strings.Contains(contentType, "text/csv") {
		t.Fatalf("unexpected content type: %q", contentType)
	}
	body := w.Body.String()
	if !strings.Contains(body, "tickets.update") || !strings.Contains(body, "status") {
		t.Fatalf("unexpected csv body=%s", body)
	}
	if svc.query.TenantID != "tenant-a" || svc.query.WorkspaceID != "workspace-1" {
		t.Fatalf("unexpected scope query: %+v", svc.query)
	}
	if svc.query.Page != 1 || svc.query.PageSize != 10 {
		t.Fatalf("unexpected paging query: %+v", svc.query)
	}
}

func TestAuditHandlerExportCSVRejectsInvalidFilters(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &stubAuditQueryService{}
	r := gin.New()
	RegisterAuditRoutes(&r.RouterGroup, NewAuditHandler(svc))

	req := httptest.NewRequest(http.MethodGet, "/audit/logs/export?success=maybe", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d body=%s", w.Code, w.Body.String())
	}
}

func TestAuditHandlerExportCSVClampsLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &stubAuditQueryService{}
	r := gin.New()
	RegisterAuditRoutes(&r.RouterGroup, NewAuditHandler(svc))

	req := httptest.NewRequest(http.MethodGet, "/audit/logs/export?limit=99999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d body=%s", w.Code, w.Body.String())
	}
	if svc.query.Page != 1 || svc.query.PageSize != 5000 {
		t.Fatalf("unexpected paging query: %+v", svc.query)
	}
}
