package infra

import (
	"context"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"servify/apps/server/internal/models"
	platformauth "servify/apps/server/internal/platform/auth"
)

func newAnalyticsScopeTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:analytics_scope?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(
		&models.User{},
		&models.Customer{},
		&models.Agent{},
		&models.Ticket{},
		&models.Session{},
		&models.Message{},
		&models.CustomerSatisfaction{},
		&models.DailyStats{},
	); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	return db
}

func scopedAnalyticsContext(tenantID, workspaceID string) context.Context {
	return platformauth.ContextWithScope(context.Background(), tenantID, workspaceID)
}

func TestGormRepositoryScopedDashboardIgnoresGlobalDailyStats(t *testing.T) {
	db := newAnalyticsScopeTestDB(t)
	repo := NewGormRepository(db)
	now := time.Now()
	today := now.Truncate(24 * time.Hour)
	ctxA := scopedAnalyticsContext("tenant-a", "workspace-a")

	customerUser := models.User{ID: 11, Username: "cust-a", Email: "cust-a@example.com", Role: "customer"}
	if err := db.Create(&customerUser).Error; err != nil {
		t.Fatalf("create customer user: %v", err)
	}
	customer := models.Customer{UserID: customerUser.ID, TenantID: "tenant-a", WorkspaceID: "workspace-a", Priority: "normal"}
	if err := db.Create(&customer).Error; err != nil {
		t.Fatalf("create customer: %v", err)
	}
	agentUser := models.User{ID: 12, Username: "agent-a", Email: "agent-a@example.com", Role: "agent"}
	if err := db.Create(&agentUser).Error; err != nil {
		t.Fatalf("create agent user: %v", err)
	}
	agent := models.Agent{UserID: agentUser.ID, TenantID: "tenant-a", WorkspaceID: "workspace-a", Status: "online", AvgResponseTime: 42}
	if err := db.Create(&agent).Error; err != nil {
		t.Fatalf("create agent: %v", err)
	}
	ticket := models.Ticket{
		ID:          101,
		TenantID:    "tenant-a",
		WorkspaceID: "workspace-a",
		Title:       "Scoped ticket",
		Status:      "open",
		CustomerID:  customerUser.ID,
		CreatedAt:   today.Add(2 * time.Hour),
		UpdatedAt:   today.Add(2 * time.Hour),
	}
	if err := db.Create(&ticket).Error; err != nil {
		t.Fatalf("create ticket: %v", err)
	}
	session := models.Session{
		ID:          "sess-a",
		TenantID:    "tenant-a",
		WorkspaceID: "workspace-a",
		UserID:      customerUser.ID,
		Status:      "active",
		CreatedAt:   today.Add(2 * time.Hour),
		StartedAt:   today.Add(2 * time.Hour),
	}
	if err := db.Create(&session).Error; err != nil {
		t.Fatalf("create session: %v", err)
	}
	if err := db.Create(&models.Message{
		TenantID:    "tenant-a",
		WorkspaceID: "workspace-a",
		SessionID:   session.ID,
		Content:     "hello",
		Type:        "text",
		Sender:      "customer",
		CreatedAt:   today.Add(2 * time.Hour),
	}).Error; err != nil {
		t.Fatalf("create message: %v", err)
	}
	// Global/system DailyStats contains unrelated high counters and must not leak
	// into scoped dashboard results.
	if err := db.Create(&models.DailyStats{
		Date:              today,
		AIUsageCount:      999,
		WeKnoraUsageCount: 555,
	}).Error; err != nil {
		t.Fatalf("create daily stats: %v", err)
	}

	stats, err := repo.GetDashboardStats(ctxA)
	if err != nil {
		t.Fatalf("GetDashboardStats: %v", err)
	}
	if stats.TotalCustomers != 1 || stats.TotalAgents != 1 || stats.TotalTickets != 1 || stats.TotalSessions != 1 {
		t.Fatalf("unexpected scoped counts: %+v", stats)
	}
	if stats.AIUsageToday != 0 || stats.WeKnoraUsageToday != 0 {
		t.Fatalf("expected scoped dashboard to ignore global daily stats, got %+v", stats)
	}
}

func TestGormRepositoryScopedTimeRangeIgnoresGlobalDailyStats(t *testing.T) {
	db := newAnalyticsScopeTestDB(t)
	repo := NewGormRepository(db)
	day := time.Date(2026, 4, 8, 0, 0, 0, 0, time.UTC)
	ctxA := scopedAnalyticsContext("tenant-a", "workspace-a")

	if err := db.Create(&models.CustomerSatisfaction{
		TenantID:    "tenant-a",
		WorkspaceID: "workspace-a",
		TicketID:    201,
		CustomerID:  21,
		Rating:      4,
		Category:    "overall",
		CreatedAt:   day.Add(4 * time.Hour),
	}).Error; err != nil {
		t.Fatalf("create scoped satisfaction: %v", err)
	}
	if err := db.Create(&models.CustomerSatisfaction{
		TenantID:    "tenant-b",
		WorkspaceID: "workspace-b",
		TicketID:    202,
		CustomerID:  22,
		Rating:      1,
		Category:    "overall",
		CreatedAt:   day.Add(5 * time.Hour),
	}).Error; err != nil {
		t.Fatalf("create unscoped satisfaction: %v", err)
	}
	if err := db.Create(&models.DailyStats{
		Date:                 day,
		AvgResponseTime:      777,
		CustomerSatisfaction: 1.5,
	}).Error; err != nil {
		t.Fatalf("create daily stats: %v", err)
	}

	items, err := repo.GetTimeRangeStats(ctxA, day, day)
	if err != nil {
		t.Fatalf("GetTimeRangeStats: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 stat row, got %d", len(items))
	}
	if items[0].CustomerSatisfaction != 4 {
		t.Fatalf("expected scoped satisfaction avg 4, got %+v", items[0])
	}
	if items[0].AvgResponseTime != 0 {
		t.Fatalf("expected scoped time-range to avoid global avg response time, got %+v", items[0])
	}
}
