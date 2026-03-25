package audit

import (
	"context"
	"testing"
	"time"

	"servify/apps/server/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func openTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&models.AuditLog{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func TestGormQueryServiceList(t *testing.T) {
	db := openTestDB(t)
	now := time.Now().UTC()
	logs := []models.AuditLog{
		{PrincipalKind: "agent", Action: "tickets.create", ResourceType: "tickets", ResourceID: "1", Success: true, CreatedAt: now.Add(-2 * time.Hour)},
		{PrincipalKind: "service", Action: "metrics_ingest.create", ResourceType: "metrics_ingest", ResourceID: "job-1", Success: true, CreatedAt: now.Add(-1 * time.Hour)},
		{PrincipalKind: "admin", Action: "tickets.assign", ResourceType: "tickets", ResourceID: "2", Success: false, CreatedAt: now},
	}
	if err := db.Create(&logs).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}

	svc := NewGormQueryService(db)
	success := true
	items, total, err := svc.List(context.Background(), ListQuery{
		ResourceType: "tickets",
		Success:      &success,
		Page:         1,
		PageSize:     10,
	})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if total != 1 {
		t.Fatalf("total = %d want 1", total)
	}
	if len(items) != 1 || items[0].Action != "tickets.create" {
		t.Fatalf("unexpected items: %+v", items)
	}
}
