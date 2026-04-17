//go:build integration
// +build integration

package services

import (
	"context"
	"testing"

	"servify/apps/server/internal/models"

	"github.com/sirupsen/logrus"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func newIntegrationTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&models.AppIntegration{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	return db
}

func TestAppIntegrationService_CreateAndList(t *testing.T) {
	db := newIntegrationTestDB(t)
	svc := NewAppIntegrationService(db, logrus.New())

	item, err := svc.Create(context.Background(), &AppIntegrationCreateRequest{
		Name:      "Stripe Billing",
		Slug:      "stripe-billing",
		Vendor:    "Stripe",
		IFrameURL: "https://example.com/stripe",
	})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if item.Slug != "stripe-billing" {
		t.Fatalf("unexpected slug: %s", item.Slug)
	}

	list, total, err := svc.List(context.Background(), &AppIntegrationListRequest{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if total != 1 || len(list) != 1 {
		t.Fatalf("unexpected list result: total=%d len=%d", total, len(list))
	}
}

func TestAppIntegrationService_ListScopedByWorkspace(t *testing.T) {
	db := newIntegrationTestDB(t)
	svc := NewAppIntegrationService(db, logrus.New())

	ctxA := scopedContext("tenant-a", "workspace-a")
	ctxB := scopedContext("tenant-a", "workspace-b")

	if _, err := svc.Create(ctxA, &AppIntegrationCreateRequest{Name: "Stripe A", Slug: "stripe-a", IFrameURL: "https://a.example.com"}); err != nil {
		t.Fatalf("create A failed: %v", err)
	}
	if _, err := svc.Create(ctxB, &AppIntegrationCreateRequest{Name: "Stripe B", Slug: "stripe-b", IFrameURL: "https://b.example.com"}); err != nil {
		t.Fatalf("create B failed: %v", err)
	}

	list, total, err := svc.List(ctxA, &AppIntegrationListRequest{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("list A failed: %v", err)
	}
	if total != 1 || len(list) != 1 || list[0].Slug != "stripe-a" {
		t.Fatalf("unexpected list result: total=%d list=%+v", total, list)
	}
}
