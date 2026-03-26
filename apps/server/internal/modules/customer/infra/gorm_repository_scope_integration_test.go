//go:build integration
// +build integration

package infra

import (
	"context"
	"strings"
	"testing"
	"time"

	"servify/apps/server/internal/models"
	customerapp "servify/apps/server/internal/modules/customer/application"
	platformauth "servify/apps/server/internal/platform/auth"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newCustomerInfraTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := "file:customer_infra_" + strings.ReplaceAll(t.Name(), "/", "_") + "?mode=memory&cache=shared"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.Customer{}, &models.Session{}, &models.Ticket{}, &models.Message{}); err != nil {
		t.Fatalf("automigrate: %v", err)
	}
	return db
}

func scopedCustomerContext(tenantID, workspaceID string) context.Context {
	return platformauth.ContextWithScope(context.Background(), tenantID, workspaceID)
}

func TestCustomerRepositoryAppliesScopeOnCreateAndRead(t *testing.T) {
	db := newCustomerInfraTestDB(t)
	repo := NewGormRepository(db)
	ctxA := scopedCustomerContext("tenant-a", "workspace-a")
	ctxB := scopedCustomerContext("tenant-b", "workspace-b")

	user, err := repo.CreateCustomer(ctxA, customerapp.CreateCustomerCommand{
		Username: "alice",
		Email:    "alice@example.com",
		Name:     "Alice",
		Source:   "web",
		Priority: "normal",
	})
	if err != nil {
		t.Fatalf("create customer: %v", err)
	}

	var stored models.Customer
	if err := db.First(&stored, "user_id = ?", user.ID).Error; err != nil {
		t.Fatalf("load customer: %v", err)
	}
	if stored.TenantID != "tenant-a" || stored.WorkspaceID != "workspace-a" {
		t.Fatalf("unexpected customer scope: %+v", stored)
	}

	if _, err := repo.GetCustomerByID(ctxA, user.ID); err != nil {
		t.Fatalf("get customer with matching scope: %v", err)
	}
	if _, err := repo.GetCustomerByID(ctxB, user.ID); err == nil {
		t.Fatal("expected cross-tenant customer lookup to fail")
	}
}

func TestCustomerRepositoryFiltersListAndStatsByScope(t *testing.T) {
	db := newCustomerInfraTestDB(t)
	repo := NewGormRepository(db)
	now := time.Now()

	users := []models.User{
		{ID: 1, Username: "alice", Email: "alice@example.com", Role: "customer", Status: "active", CreatedAt: now, UpdatedAt: now},
		{ID: 2, Username: "bob", Email: "bob@example.com", Role: "customer", Status: "active", CreatedAt: now, UpdatedAt: now},
	}
	if err := db.Create(&users).Error; err != nil {
		t.Fatalf("seed users: %v", err)
	}
	customers := []models.Customer{
		{UserID: 1, TenantID: "tenant-a", WorkspaceID: "workspace-a", Company: "A Co", Source: "web", Priority: "normal", CreatedAt: now, UpdatedAt: now},
		{UserID: 2, TenantID: "tenant-b", WorkspaceID: "workspace-b", Company: "B Co", Source: "referral", Priority: "high", CreatedAt: now, UpdatedAt: now},
	}
	if err := db.Create(&customers).Error; err != nil {
		t.Fatalf("seed customers: %v", err)
	}

	ctxA := scopedCustomerContext("tenant-a", "workspace-a")
	items, total, err := repo.ListCustomers(ctxA, customerapp.ListCustomersQuery{
		Page:      1,
		PageSize:  20,
		SortBy:    "created_at",
		SortOrder: "desc",
	})
	if err != nil {
		t.Fatalf("list customers: %v", err)
	}
	if total != 1 || len(items) != 1 || items[0].User.ID != 1 {
		t.Fatalf("unexpected scoped customer list: total=%d items=%+v", total, items)
	}

	stats, err := repo.GetStats(ctxA)
	if err != nil {
		t.Fatalf("get stats: %v", err)
	}
	if stats.Total != 1 {
		t.Fatalf("unexpected scoped stats: %+v", stats)
	}
}
