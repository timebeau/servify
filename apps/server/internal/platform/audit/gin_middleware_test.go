package audit

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

type stubRecorder struct {
	entries []Entry
}

func (s *stubRecorder) Record(_ context.Context, entry Entry) error {
	s.entries = append(s.entries, entry)
	return nil
}

func TestMiddlewareRecordsSuccessfulWrite(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := &stubRecorder{}

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uint(7))
		c.Set("principal_kind", "agent")
		c.Next()
	})
	r.Use(Middleware(recorder))
	r.POST("/api/tickets/:id/assign", func(c *gin.Context) {
		SetBefore(c, gin.H{"agent_id": nil})
		SetAfter(c, gin.H{"agent_id": 42})
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodPost, "/api/tickets/123/assign", strings.NewReader(`{"target_agent_id":42}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-ID", "req-1")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d body=%s", w.Code, w.Body.String())
	}
	if len(recorder.entries) != 1 {
		t.Fatalf("expected 1 audit entry got %d", len(recorder.entries))
	}
	entry := recorder.entries[0]
	if entry.ResourceType != "tickets" {
		t.Fatalf("resource_type = %q want tickets", entry.ResourceType)
	}
	if entry.ResourceID != "123" {
		t.Fatalf("resource_id = %q want 123", entry.ResourceID)
	}
	if entry.Action != "tickets.assign" {
		t.Fatalf("action = %q want tickets.assign", entry.Action)
	}
	if entry.PrincipalKind != "agent" {
		t.Fatalf("principal_kind = %q want agent", entry.PrincipalKind)
	}
	if entry.ActorUserID == nil || *entry.ActorUserID != 7 {
		t.Fatalf("actor_user_id = %v want 7", entry.ActorUserID)
	}
	if entry.RequestID != "req-1" {
		t.Fatalf("request_id = %q want req-1", entry.RequestID)
	}
	if entry.RequestJSON != `{"target_agent_id":42}` {
		t.Fatalf("request_json = %q", entry.RequestJSON)
	}
	if entry.BeforeJSON == "" || entry.AfterJSON == "" {
		t.Fatalf("expected before/after json, got before=%q after=%q", entry.BeforeJSON, entry.AfterJSON)
	}
}

func TestMiddlewareSkipsFailedWriteAndReads(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := &stubRecorder{}

	r := gin.New()
	r.Use(Middleware(recorder))
	r.GET("/api/tickets", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	r.POST("/api/tickets", func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
	})

	wRead := httptest.NewRecorder()
	r.ServeHTTP(wRead, httptest.NewRequest(http.MethodGet, "/api/tickets", nil))

	wWrite := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/tickets", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(wWrite, req)

	if len(recorder.entries) != 0 {
		t.Fatalf("expected no audit entries, got %d", len(recorder.entries))
	}
}

func TestMiddlewareRedactsSensitiveFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := &stubRecorder{}

	r := gin.New()
	r.Use(Middleware(recorder))
	r.POST("/api/apps/config", func(c *gin.Context) {
		SetBefore(c, gin.H{"api_key": "before-secret"})
		SetAfter(c, gin.H{"nested": gin.H{"access_token": "after-secret"}})
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodPost, "/api/apps/config", strings.NewReader(`{"password":"p@ss","profile":{"api_key":"k-1"},"list":[{"refresh_token":"r-1"}]}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d body=%s", w.Code, w.Body.String())
	}
	if len(recorder.entries) != 1 {
		t.Fatalf("expected 1 audit entry got %d", len(recorder.entries))
	}

	entry := recorder.entries[0]
	if strings.Contains(entry.RequestJSON, "p@ss") || strings.Contains(entry.RequestJSON, "k-1") || strings.Contains(entry.RequestJSON, "r-1") {
		t.Fatalf("request_json should redact secrets, got %q", entry.RequestJSON)
	}
	if strings.Contains(entry.BeforeJSON, "before-secret") {
		t.Fatalf("before_json should redact secrets, got %q", entry.BeforeJSON)
	}
	if strings.Contains(entry.AfterJSON, "after-secret") {
		t.Fatalf("after_json should redact secrets, got %q", entry.AfterJSON)
	}
	if !strings.Contains(entry.RequestJSON, "[REDACTED]") || !strings.Contains(entry.BeforeJSON, "[REDACTED]") || !strings.Contains(entry.AfterJSON, "[REDACTED]") {
		t.Fatalf("expected redaction markers, got request=%q before=%q after=%q", entry.RequestJSON, entry.BeforeJSON, entry.AfterJSON)
	}
}

func TestMiddlewareAllowsAuditOverrides(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := &stubRecorder{}

	r := gin.New()
	r.Use(Middleware(recorder))
	r.PUT("/security/config/tenant", func(c *gin.Context) {
		SetAction(c, "scoped_config.tenant.update")
		SetResourceType(c, "scoped_config")
		SetResourceID(c, "tenant-a")
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/security/config/tenant", strings.NewReader(`{"portal":{"brand_name":"Tenant"}}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d body=%s", w.Code, w.Body.String())
	}
	if len(recorder.entries) != 1 {
		t.Fatalf("expected 1 audit entry got %d", len(recorder.entries))
	}
	entry := recorder.entries[0]
	if entry.Action != "scoped_config.tenant.update" {
		t.Fatalf("action = %q want scoped_config.tenant.update", entry.Action)
	}
	if entry.ResourceType != "scoped_config" {
		t.Fatalf("resource_type = %q want scoped_config", entry.ResourceType)
	}
	if entry.ResourceID != "tenant-a" {
		t.Fatalf("resource_id = %q want tenant-a", entry.ResourceID)
	}
}
