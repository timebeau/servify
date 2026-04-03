//go:build integration
// +build integration

package services

import (
	"context"
	"testing"
	"time"

	"servify/apps/server/internal/models"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newSLATestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&models.Ticket{}, &models.SLAConfig{}, &models.SLAViolation{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	return db
}

func TestSLAService_CheckViolation_FirstResponse(t *testing.T) {
	db := newSLATestDB(t)
	svc := NewSLAService(db, logrus.New())

	now := time.Now()
	cfg := &models.SLAConfig{
		Name:              "High Priority",
		Priority:          "high",
		FirstResponseTime: 5,
		ResolutionTime:    60,
		EscalationTime:    30,
		Active:            true,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	if err := db.Create(cfg).Error; err != nil {
		t.Fatalf("failed to insert config: %v", err)
	}

	ticket := &models.Ticket{
		ID:        1,
		Title:     "Test",
		Priority:  "high",
		Status:    "open",
		CreatedAt: now.Add(-10 * time.Minute),
		UpdatedAt: now.Add(-10 * time.Minute),
	}
	if err := db.Create(ticket).Error; err != nil {
		t.Fatalf("failed to insert ticket: %v", err)
	}

	violation, err := svc.CheckSLAViolation(context.Background(), ticket)
	if err != nil {
		t.Fatalf("CheckSLAViolation failed: %v", err)
	}
	if violation == nil || violation.ViolationType != "first_response" {
		t.Fatalf("expected first_response violation, got %+v", violation)
	}

	// running again should reuse existing violation
	violation2, err := svc.CheckSLAViolation(context.Background(), ticket)
	if err != nil {
		t.Fatalf("CheckSLAViolation second call failed: %v", err)
	}
	if violation2 == nil || violation2.ID != violation.ID {
		t.Fatalf("expected duplicate detection; got %#v", violation2)
	}
}

func TestSLAService_ListConfigsScopedByWorkspace(t *testing.T) {
	db := newSLATestDB(t)
	svc := NewSLAService(db, logrus.New())

	ctxA := scopedContext("tenant-a", "workspace-a")
	ctxB := scopedContext("tenant-a", "workspace-b")

	if _, err := svc.CreateSLAConfig(ctxA, &SLAConfigCreateRequest{
		Name: "A", Priority: "high", FirstResponseTime: 5, ResolutionTime: 60, EscalationTime: 30,
	}); err != nil {
		t.Fatalf("create A failed: %v", err)
	}
	if _, err := svc.CreateSLAConfig(ctxB, &SLAConfigCreateRequest{
		Name: "B", Priority: "normal", FirstResponseTime: 10, ResolutionTime: 120, EscalationTime: 30,
	}); err != nil {
		t.Fatalf("create B failed: %v", err)
	}

	items, total, err := svc.ListSLAConfigs(ctxA, &SLAConfigListRequest{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if total != 1 || len(items) != 1 || items[0].WorkspaceID != "workspace-a" {
		t.Fatalf("unexpected items: total=%d items=%+v", total, items)
	}
}

func TestSLAService_ResolveViolationsByTicket(t *testing.T) {
	db := newSLATestDB(t)
	svc := NewSLAService(db, logrus.New())

	now := time.Now()
	cfg := &models.SLAConfig{
		Name:              "Normal Priority",
		Priority:          "normal",
		FirstResponseTime: 10,
		ResolutionTime:    120,
		EscalationTime:    60,
		Active:            true,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	if err := db.Create(cfg).Error; err != nil {
		t.Fatalf("failed to insert config: %v", err)
	}

	ticket := &models.Ticket{
		ID:        2,
		Title:     "Need help",
		Priority:  "normal",
		Status:    "open",
		CreatedAt: now.Add(-3 * time.Hour),
		UpdatedAt: now.Add(-3 * time.Hour),
	}
	if err := db.Create(ticket).Error; err != nil {
		t.Fatalf("failed to insert ticket: %v", err)
	}

	// create violation manually via service
	if _, err := svc.CheckSLAViolation(context.Background(), ticket); err != nil {
		t.Fatalf("failed to create violation: %v", err)
	}

	if err := svc.ResolveViolationsByTicket(context.Background(), ticket.ID, []string{"first_response"}); err != nil {
		t.Fatalf("ResolveViolationsByTicket failed: %v", err)
	}

	var count int64
	if err := db.Model(&models.SLAViolation{}).Where("ticket_id = ? AND resolved = true", ticket.ID).Count(&count).Error; err != nil {
		t.Fatalf("failed to count resolved violations: %v", err)
	}
	if count == 0 {
		t.Fatalf("expected resolved violations, got zero")
	}
}

func TestSLAService_ViolationsScopedByWorkspace(t *testing.T) {
	db := newSLATestDB(t)
	svc := NewSLAService(db, logrus.New())

	now := time.Now()
	ctxA := scopedContext("tenant-a", "workspace-a")
	ctxB := scopedContext("tenant-a", "workspace-b")

	cfgA, err := svc.CreateSLAConfig(ctxA, &SLAConfigCreateRequest{
		Name: "A", Priority: "high", FirstResponseTime: 5, ResolutionTime: 60, EscalationTime: 30,
	})
	if err != nil {
		t.Fatalf("create config A failed: %v", err)
	}
	cfgB, err := svc.CreateSLAConfig(ctxB, &SLAConfigCreateRequest{
		Name: "B", Priority: "urgent", FirstResponseTime: 5, ResolutionTime: 60, EscalationTime: 30,
	})
	if err != nil {
		t.Fatalf("create config B failed: %v", err)
	}

	ticketA := &models.Ticket{
		ID:          21,
		Title:       "Ticket A",
		Priority:    "high",
		Status:      "open",
		TenantID:    "tenant-a",
		WorkspaceID: "workspace-a",
		CreatedAt:   now.Add(-2 * time.Hour),
		UpdatedAt:   now.Add(-2 * time.Hour),
	}
	ticketB := &models.Ticket{
		ID:          22,
		Title:       "Ticket B",
		Priority:    "urgent",
		Status:      "open",
		TenantID:    "tenant-a",
		WorkspaceID: "workspace-b",
		CreatedAt:   now.Add(-2 * time.Hour),
		UpdatedAt:   now.Add(-2 * time.Hour),
	}
	if err := db.Create(ticketA).Error; err != nil {
		t.Fatalf("create ticket A failed: %v", err)
	}
	if err := db.Create(ticketB).Error; err != nil {
		t.Fatalf("create ticket B failed: %v", err)
	}

	violationA, err := svc.CheckSLAViolation(ctxA, ticketA)
	if err != nil {
		t.Fatalf("check violation A failed: %v", err)
	}
	if violationA == nil || violationA.WorkspaceID != "workspace-a" || violationA.SLAConfigID != cfgA.ID {
		t.Fatalf("unexpected violation A: %+v", violationA)
	}
	violationB, err := svc.CheckSLAViolation(ctxB, ticketB)
	if err != nil {
		t.Fatalf("check violation B failed: %v", err)
	}
	if violationB == nil || violationB.WorkspaceID != "workspace-b" || violationB.SLAConfigID != cfgB.ID {
		t.Fatalf("unexpected violation B: %+v", violationB)
	}

	itemsA, totalA, err := svc.ListSLAViolations(ctxA, &SLAViolationListRequest{Page: 1, PageSize: 20})
	if err != nil {
		t.Fatalf("list violations A failed: %v", err)
	}
	if totalA != 1 || len(itemsA) != 1 || itemsA[0].WorkspaceID != "workspace-a" {
		t.Fatalf("unexpected scoped violations: total=%d items=%+v", totalA, itemsA)
	}

	statsA, err := svc.GetSLAStats(ctxA)
	if err != nil {
		t.Fatalf("stats A failed: %v", err)
	}
	if statsA.TotalConfigs != 1 || statsA.TotalViolations != 1 || statsA.ViolationsByPriority["high"] != 1 {
		t.Fatalf("unexpected scoped stats: %+v", statsA)
	}

	if config, err := svc.GetSLAConfigByPriority(ctxA, "urgent", ""); err != nil {
		t.Fatalf("get urgent config in scope A errored: %v", err)
	} else if config != nil {
		t.Fatalf("expected urgent config to be hidden from workspace A, got %+v", config)
	}

	if err := svc.ResolveSLAViolation(ctxB, violationA.ID); err == nil {
		t.Fatal("expected cross-workspace resolve to fail")
	}
	if err := svc.DeleteSLAConfig(ctxB, cfgA.ID); err == nil {
		t.Fatal("expected cross-workspace delete config to fail")
	}

	if err := svc.ResolveSLAViolation(ctxA, violationA.ID); err != nil {
		t.Fatalf("resolve violation A failed: %v", err)
	}
	resolved := true
	resolvedItems, resolvedTotal, err := svc.ListSLAViolations(ctxA, &SLAViolationListRequest{
		Page: 1, PageSize: 20, Resolved: &resolved,
	})
	if err != nil {
		t.Fatalf("list resolved violations failed: %v", err)
	}
	if resolvedTotal != 1 || len(resolvedItems) != 1 || !resolvedItems[0].Resolved {
		t.Fatalf("unexpected resolved items: total=%d items=%+v", resolvedTotal, resolvedItems)
	}
}
