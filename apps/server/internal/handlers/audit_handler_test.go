package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"servify/apps/server/internal/models"
	auditplatform "servify/apps/server/internal/platform/audit"

	"github.com/gin-gonic/gin"
)

type stubAuditQueryService struct {
	items []models.AuditLog
	total int64
	err   error
	query auditplatform.ListQuery
}

func (s *stubAuditQueryService) List(_ context.Context, query auditplatform.ListQuery) ([]models.AuditLog, int64, error) {
	s.query = query
	return s.items, s.total, s.err
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
