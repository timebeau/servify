package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"servify/apps/server/internal/models"
	conversationdelivery "servify/apps/server/internal/modules/conversation/delivery"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func newSessionTransferLegacyStoreTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&models.Session{}, &models.Ticket{}, &models.TicketStatus{}, &models.Agent{}); err != nil {
		t.Fatalf("automigrate: %v", err)
	}
	return db
}

func TestSessionTransferLegacyStoreSyncTransferSessionReturnsNotFound(t *testing.T) {
	db := newSessionTransferLegacyStoreTestDB(t)
	store := newSessionTransferLegacyStore(db)

	err := db.Transaction(func(tx *gorm.DB) error {
		return store.syncTransferSession(context.Background(), tx, &conversationdelivery.TransferSession{ID: "missing"}, 7)
	})
	if err == nil || !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected ErrRecordNotFound, got %v", err)
	}
}

func TestSessionTransferLegacyStoreSyncTransferTicketReturnsNotFound(t *testing.T) {
	db := newSessionTransferLegacyStoreTestDB(t)
	store := newSessionTransferLegacyStore(db)

	err := db.Transaction(func(tx *gorm.DB) error {
		return store.syncTransferTicket(tx, 999, 7, time.Now())
	})
	if err == nil || !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected ErrRecordNotFound, got %v", err)
	}
}
