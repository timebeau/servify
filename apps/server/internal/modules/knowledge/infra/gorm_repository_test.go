//go:build integration
// +build integration

package infra

import (
	"context"
	"testing"
	"time"

	"servify/apps/server/internal/models"
	"servify/apps/server/internal/modules/knowledge/domain"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newKnowledgeInfraTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&models.KnowledgeDoc{}, &models.KnowledgeIndexJob{}); err != nil {
		t.Fatalf("automigrate: %v", err)
	}
	return db
}

func TestGormIndexJobRepositoryRoundTrip(t *testing.T) {
	db := newKnowledgeInfraTestDB(t)
	docRepo := NewGormDocumentRepository(db)
	jobRepo := NewGormIndexJobRepository(db)

	doc := &domain.Document{
		Title:     "Billing",
		Content:   "Billing details",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := docRepo.Create(context.Background(), doc); err != nil {
		t.Fatalf("create doc: %v", err)
	}

	job := &domain.IndexJob{
		ID:         "job-1",
		DocumentID: doc.ID,
		Status:     domain.IndexJobQueued,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	if err := jobRepo.Create(context.Background(), job); err != nil {
		t.Fatalf("create job: %v", err)
	}

	completedAt := time.Now()
	job.Status = domain.IndexJobDone
	job.CompletedAt = &completedAt
	job.UpdatedAt = completedAt
	if err := jobRepo.Update(context.Background(), job); err != nil {
		t.Fatalf("update job: %v", err)
	}

	got, err := jobRepo.Get(context.Background(), job.ID)
	if err != nil {
		t.Fatalf("get job: %v", err)
	}
	if got.DocumentID != doc.ID {
		t.Fatalf("document id = %s want %s", got.DocumentID, doc.ID)
	}
	if got.Status != domain.IndexJobDone {
		t.Fatalf("status = %s want %s", got.Status, domain.IndexJobDone)
	}
	if got.CompletedAt == nil {
		t.Fatal("expected completed_at to be persisted")
	}
}
