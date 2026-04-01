package audit

import (
	"context"
	"testing"
	"time"

	"servify/apps/server/internal/models"
)

func TestGormRetentionServiceCleanup(t *testing.T) {
	db := openTestDB(t)
	now := time.Now().UTC()
	logs := []models.AuditLog{
		{Action: "old-1", PrincipalKind: "admin", ResourceType: "ticket", Route: "/api/tickets/1", Method: "POST", CreatedAt: now.Add(-400 * 24 * time.Hour)},
		{Action: "old-2", PrincipalKind: "admin", ResourceType: "ticket", Route: "/api/tickets/2", Method: "POST", CreatedAt: now.Add(-300 * 24 * time.Hour)},
		{Action: "fresh", PrincipalKind: "admin", ResourceType: "ticket", Route: "/api/tickets/3", Method: "POST", CreatedAt: now.Add(-30 * 24 * time.Hour)},
	}
	if err := db.Create(&logs).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}

	svc := NewGormRetentionService(db, 180*24*time.Hour, 1)
	deleted, err := svc.Cleanup(context.Background(), now)
	if err != nil {
		t.Fatalf("Cleanup() error = %v", err)
	}
	if deleted != 2 {
		t.Fatalf("deleted = %d want 2", deleted)
	}

	var remaining []models.AuditLog
	if err := db.Order("created_at asc").Find(&remaining).Error; err != nil {
		t.Fatalf("query remaining: %v", err)
	}
	if len(remaining) != 1 || remaining[0].Action != "fresh" {
		t.Fatalf("unexpected remaining logs: %+v", remaining)
	}
}

func TestGormRetentionServiceCleanupNoopWhenDisabled(t *testing.T) {
	db := openTestDB(t)
	svc := NewGormRetentionService(db, 0, 10)
	if svc != nil {
		t.Fatal("expected nil service for non-positive retention")
	}
}
