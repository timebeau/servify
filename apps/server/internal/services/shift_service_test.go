//go:build integration
// +build integration

package services

import (
	"context"
	"testing"
	"time"

	"servify/apps/server/internal/models"

	"github.com/sirupsen/logrus"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func newTestDB(t *testing.T, tables ...interface{}) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}
	if len(tables) == 0 {
		tables = []interface{}{&models.ShiftSchedule{}, &models.User{}, &models.Agent{}}
	}
	tables = append(tables, &models.User{}, &models.Agent{})
	if err := db.AutoMigrate(tables...); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	return db
}

func TestShiftService_Create_and_GetStats(t *testing.T) {
	db := newTestDB(t, &models.ShiftSchedule{}, &models.User{}, &models.Agent{})
	svc := NewShiftService(db, logrus.New())
	if err := db.Create(&models.User{ID: 1, Username: "agent1", Email: "agent1@example.com", Role: "agent"}).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	if err := db.Create(&models.Agent{UserID: 1, Status: "online"}).Error; err != nil {
		t.Fatalf("create agent: %v", err)
	}

	now := time.Now()
	req := &ShiftCreateRequest{
		AgentID:   1,
		ShiftType: "morning",
		StartTime: now.Add(time.Hour),
		EndTime:   now.Add(4 * time.Hour),
		Status:    "scheduled",
	}
	_, err := svc.CreateShift(context.Background(), req)
	if err != nil {
		t.Fatalf("CreateShift failed: %v", err)
	}

	stats, err := svc.GetShiftStats(context.Background())
	if err != nil {
		t.Fatalf("GetShiftStats failed: %v", err)
	}
	if stats.Total != 1 {
		t.Fatalf("expected total 1, got %d", stats.Total)
	}
	if stats.ByType["morning"] != 1 {
		t.Fatalf("expected morning count 1, got %d", stats.ByType["morning"])
	}
	if stats.ByStatus["scheduled"] != 1 {
		t.Fatalf("expected scheduled count 1, got %d", stats.ByStatus["scheduled"])
	}
}

func TestShiftService_Create_Validation(t *testing.T) {
	db := newTestDB(t, &models.ShiftSchedule{}, &models.User{}, &models.Agent{})
	svc := NewShiftService(db, logrus.New())
	if err := db.Create(&models.User{ID: 1, Username: "agentv", Email: "agentv@example.com", Role: "agent"}).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	if err := db.Create(&models.Agent{UserID: 1, Status: "online"}).Error; err != nil {
		t.Fatalf("create agent: %v", err)
	}

	now := time.Now()
	req := &ShiftCreateRequest{
		AgentID:   1,
		ShiftType: "morning",
		StartTime: now,
		EndTime:   now.Add(-time.Hour), // invalid: end before start
		Status:    "scheduled",
	}
	if _, err := svc.CreateShift(context.Background(), req); err == nil {
		t.Fatalf("expected validation error for end before start")
	}
}

func TestShiftService_Update_and_Delete(t *testing.T) {
	db := newTestDB(t, &models.ShiftSchedule{}, &models.User{}, &models.Agent{})
	svc := NewShiftService(db, logrus.New())
	if err := db.Create(&models.User{ID: 2, Username: "agent2", Email: "agent2@example.com", Role: "agent"}).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	if err := db.Create(&models.Agent{UserID: 2, Status: "online"}).Error; err != nil {
		t.Fatalf("create agent: %v", err)
	}

	now := time.Now()
	shift, err := svc.CreateShift(context.Background(), &ShiftCreateRequest{
		AgentID:   2,
		ShiftType: "night",
		StartTime: now,
		EndTime:   now.Add(3 * time.Hour),
		Status:    "scheduled",
	})
	if err != nil {
		t.Fatalf("CreateShift failed: %v", err)
	}

	newStatus := "active"
	updated, err := svc.UpdateShift(context.Background(), shift.ID, &ShiftUpdateRequest{
		Status: &newStatus,
	})
	if err != nil {
		t.Fatalf("UpdateShift failed: %v", err)
	}
	if updated.Status != newStatus {
		t.Fatalf("expected status %s, got %s", newStatus, updated.Status)
	}

	if err := svc.DeleteShift(context.Background(), shift.ID); err != nil {
		t.Fatalf("DeleteShift failed: %v", err)
	}
	if err := svc.DeleteShift(context.Background(), shift.ID); err == nil {
		t.Fatalf("expected error deleting non-existent shift")
	}
}

func TestShiftService_ScopedByWorkspace(t *testing.T) {
	db := newTestDB(t, &models.ShiftSchedule{}, &models.User{}, &models.Agent{})
	svc := NewShiftService(db, logrus.New())

	now := time.Now()
	ctxA := scopedContext("tenant-a", "workspace-a")
	ctxB := scopedContext("tenant-a", "workspace-b")
	if err := db.Create(&models.User{ID: 1, Username: "agenta", Email: "agenta@example.com", Role: "agent"}).Error; err != nil {
		t.Fatalf("create user A: %v", err)
	}
	if err := db.Create(&models.User{ID: 2, Username: "agentb", Email: "agentb@example.com", Role: "agent"}).Error; err != nil {
		t.Fatalf("create user B: %v", err)
	}
	if err := db.Create(&models.Agent{TenantID: "tenant-a", WorkspaceID: "workspace-a", UserID: 1, Status: "online"}).Error; err != nil {
		t.Fatalf("create agent A: %v", err)
	}
	if err := db.Create(&models.Agent{TenantID: "tenant-a", WorkspaceID: "workspace-b", UserID: 2, Status: "online"}).Error; err != nil {
		t.Fatalf("create agent B: %v", err)
	}

	shiftA, err := svc.CreateShift(ctxA, &ShiftCreateRequest{
		AgentID:   1,
		ShiftType: "morning",
		StartTime: now.Add(time.Hour),
		EndTime:   now.Add(2 * time.Hour),
		Status:    "scheduled",
	})
	if err != nil {
		t.Fatalf("create A failed: %v", err)
	}
	if shiftA.TenantID != "tenant-a" || shiftA.WorkspaceID != "workspace-a" {
		t.Fatalf("unexpected scope on create: %+v", shiftA)
	}

	if _, err := svc.CreateShift(ctxB, &ShiftCreateRequest{
		AgentID:   2,
		ShiftType: "night",
		StartTime: now.Add(3 * time.Hour),
		EndTime:   now.Add(5 * time.Hour),
		Status:    "scheduled",
	}); err != nil {
		t.Fatalf("create B failed: %v", err)
	}

	items, total, err := svc.ListShifts(ctxA, &ShiftListRequest{Page: 1, PageSize: 20})
	if err != nil {
		t.Fatalf("list scoped shifts failed: %v", err)
	}
	if total != 1 || len(items) != 1 || items[0].WorkspaceID != "workspace-a" {
		t.Fatalf("unexpected scoped shifts: total=%d items=%+v", total, items)
	}

	stats, err := svc.GetShiftStats(ctxA)
	if err != nil {
		t.Fatalf("get scoped stats failed: %v", err)
	}
	if stats.Total != 1 || stats.ByType["morning"] != 1 {
		t.Fatalf("unexpected scoped stats: %+v", stats)
	}

	newStatus := "active"
	if _, err := svc.UpdateShift(ctxB, shiftA.ID, &ShiftUpdateRequest{Status: &newStatus}); err == nil {
		t.Fatal("expected scoped update to reject cross-workspace shift")
	}
	if err := svc.DeleteShift(ctxB, shiftA.ID); err == nil {
		t.Fatal("expected scoped delete to reject cross-workspace shift")
	}
}

func TestShiftService_ListShifts_DoesNotLeakCrossScopeAgentPreload(t *testing.T) {
	db := newTestDB(t, &models.ShiftSchedule{}, &models.User{}, &models.Agent{})
	svc := NewShiftService(db, logrus.New())
	now := time.Now()

	agentUserB := &models.User{ID: 401, Username: "agent-b", Email: "agent-b@example.com", Name: "Agent B", Role: "agent"}
	if err := db.Create(agentUserB).Error; err != nil {
		t.Fatalf("create agent user B: %v", err)
	}
	if err := db.Create(&models.Agent{
		TenantID:    "tenant-b",
		WorkspaceID: "workspace-b",
		UserID:      agentUserB.ID,
		Status:      "online",
	}).Error; err != nil {
		t.Fatalf("create agent B: %v", err)
	}
	if err := db.Create(&models.ShiftSchedule{
		TenantID:    "tenant-a",
		WorkspaceID: "workspace-a",
		AgentID:     agentUserB.ID,
		ShiftType:   "morning",
		StartTime:   now.Add(time.Hour),
		EndTime:     now.Add(2 * time.Hour),
		Date:        now.Truncate(24 * time.Hour),
		Status:      "scheduled",
	}).Error; err != nil {
		t.Fatalf("create shift A: %v", err)
	}

	items, total, err := svc.ListShifts(scopedContext("tenant-a", "workspace-a"), &ShiftListRequest{Page: 1, PageSize: 20})
	if err != nil {
		t.Fatalf("ListShifts failed: %v", err)
	}
	if total != 1 || len(items) != 1 {
		t.Fatalf("unexpected scoped shifts: total=%d items=%+v", total, items)
	}
	if items[0].Agent.ID != 0 {
		t.Fatalf("expected agent preload to stay scoped, got %+v", items[0].Agent)
	}
}

func TestShiftService_CreateShift_RejectsCrossScopeAgent(t *testing.T) {
	db := newTestDB(t, &models.ShiftSchedule{}, &models.User{}, &models.Agent{})
	svc := NewShiftService(db, logrus.New())
	now := time.Now()

	agentUserB := &models.User{ID: 501, Username: "agent-b2", Email: "agent-b2@example.com", Name: "Agent B2", Role: "agent"}
	if err := db.Create(agentUserB).Error; err != nil {
		t.Fatalf("create agent user B2: %v", err)
	}
	if err := db.Create(&models.Agent{
		TenantID:    "tenant-b",
		WorkspaceID: "workspace-b",
		UserID:      agentUserB.ID,
		Status:      "online",
	}).Error; err != nil {
		t.Fatalf("create agent B2: %v", err)
	}

	if _, err := svc.CreateShift(scopedContext("tenant-a", "workspace-a"), &ShiftCreateRequest{
		AgentID:   agentUserB.ID,
		ShiftType: "morning",
		StartTime: now.Add(time.Hour),
		EndTime:   now.Add(2 * time.Hour),
		Status:    "scheduled",
	}); err == nil {
		t.Fatal("expected cross-scope agent assignment to fail")
	}
}
