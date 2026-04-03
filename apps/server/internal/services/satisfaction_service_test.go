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

func newSatisfactionTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.Customer{}, &models.Agent{}, &models.Ticket{}, &models.CustomerSatisfaction{}, &models.SatisfactionSurvey{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	return db
}

func TestSatisfactionService_ScheduleAndRespondSurvey(t *testing.T) {
	db := newSatisfactionTestDB(t)
	logger := logrus.New()
	svc := NewSatisfactionService(db, logger)

	now := time.Now()
	user := &models.User{ID: 1, Username: "customer1", Email: "c1@example.com"}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	if err := db.Create(&models.Customer{UserID: user.ID}).Error; err != nil {
		t.Fatalf("create customer: %v", err)
	}

	ticket := &models.Ticket{
		ID:         1,
		Title:      "测试工单",
		CustomerID: user.ID,
		Status:     "closed",
		Source:     "web",
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := db.Create(ticket).Error; err != nil {
		t.Fatalf("create ticket: %v", err)
	}

	survey, err := svc.ScheduleSurvey(context.Background(), ticket)
	if err != nil {
		t.Fatalf("ScheduleSurvey failed: %v", err)
	}
	if survey == nil || survey.SurveyToken == "" {
		t.Fatalf("expected survey with token, got %+v", survey)
	}
	if survey.Status != "sent" {
		t.Fatalf("expected status sent, got %s", survey.Status)
	}

	satisfaction, err := svc.RespondSurvey(context.Background(), survey.SurveyToken, 5, "great job")
	if err != nil {
		t.Fatalf("RespondSurvey failed: %v", err)
	}
	if satisfaction == nil || satisfaction.Rating != 5 {
		t.Fatalf("unexpected satisfaction result: %+v", satisfaction)
	}

	var savedSurvey models.SatisfactionSurvey
	if err := db.First(&savedSurvey, survey.ID).Error; err != nil {
		t.Fatalf("load saved survey: %v", err)
	}
	if savedSurvey.Status != "completed" {
		t.Fatalf("expected completed status, got %s", savedSurvey.Status)
	}
	if savedSurvey.SatisfactionID == nil {
		t.Fatalf("expected satisfaction_id recorded")
	}
}

func TestSatisfactionService_GetSurveyPreview(t *testing.T) {
	db := newSatisfactionTestDB(t)
	logger := logrus.New()
	svc := NewSatisfactionService(db, logger)

	now := time.Now()
	user := &models.User{ID: 2, Username: "customer2", Email: "c2@example.com"}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	if err := db.Create(&models.Customer{UserID: user.ID}).Error; err != nil {
		t.Fatalf("create customer: %v", err)
	}

	ticket := &models.Ticket{
		ID:         2,
		Title:      "预览工单",
		CustomerID: user.ID,
		Status:     "closed",
		Source:     "email",
		ResolvedAt: &now,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := db.Create(ticket).Error; err != nil {
		t.Fatalf("create ticket: %v", err)
	}

	survey, err := svc.ScheduleSurvey(context.Background(), ticket)
	if err != nil {
		t.Fatalf("ScheduleSurvey failed: %v", err)
	}

	preview, err := svc.GetSurveyPreviewByToken(context.Background(), survey.SurveyToken)
	if err != nil {
		t.Fatalf("GetSurveyPreviewByToken failed: %v", err)
	}
	if preview == nil || preview.TicketTitle != ticket.Title {
		t.Fatalf("unexpected preview result: %+v", preview)
	}
	if preview.Status != "sent" {
		t.Fatalf("expected sent status, got %s", preview.Status)
	}
	if preview.ResolvedAt == nil {
		t.Fatalf("expected resolved time in preview")
	}
}

func TestSatisfactionService_ScopedByWorkspace(t *testing.T) {
	db := newSatisfactionTestDB(t)
	logger := logrus.New()
	svc := NewSatisfactionService(db, logger)

	now := time.Now()
	customerA := &models.User{ID: 10, Username: "customer-a", Email: "customer-a@example.com"}
	customerB := &models.User{ID: 11, Username: "customer-b", Email: "customer-b@example.com"}
	if err := db.Create(customerA).Error; err != nil {
		t.Fatalf("create customer A: %v", err)
	}
	if err := db.Create(customerB).Error; err != nil {
		t.Fatalf("create customer B: %v", err)
	}
	if err := db.Create(&models.Customer{UserID: customerA.ID, TenantID: "tenant-a", WorkspaceID: "workspace-a"}).Error; err != nil {
		t.Fatalf("create customer profile A: %v", err)
	}
	if err := db.Create(&models.Customer{UserID: customerB.ID, TenantID: "tenant-a", WorkspaceID: "workspace-b"}).Error; err != nil {
		t.Fatalf("create customer profile B: %v", err)
	}

	ticketA := &models.Ticket{
		ID:          10,
		Title:       "工单 A",
		CustomerID:  customerA.ID,
		Status:      "closed",
		Source:      "web",
		TenantID:    "tenant-a",
		WorkspaceID: "workspace-a",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	ticketB := &models.Ticket{
		ID:          11,
		Title:       "工单 B",
		CustomerID:  customerB.ID,
		Status:      "closed",
		Source:      "web",
		TenantID:    "tenant-a",
		WorkspaceID: "workspace-b",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := db.Create(ticketA).Error; err != nil {
		t.Fatalf("create ticket A: %v", err)
	}
	if err := db.Create(ticketB).Error; err != nil {
		t.Fatalf("create ticket B: %v", err)
	}

	ctxA := scopedContext("tenant-a", "workspace-a")
	ctxB := scopedContext("tenant-a", "workspace-b")

	surveyA, err := svc.ScheduleSurvey(ctxA, ticketA)
	if err != nil {
		t.Fatalf("schedule survey A: %v", err)
	}
	if surveyA.TenantID != "tenant-a" || surveyA.WorkspaceID != "workspace-a" {
		t.Fatalf("unexpected scope on survey A: %+v", surveyA)
	}
	if _, err := svc.ScheduleSurvey(ctxB, ticketB); err != nil {
		t.Fatalf("schedule survey B: %v", err)
	}

	satisfactionA, err := svc.RespondSurvey(context.Background(), surveyA.SurveyToken, 5, "great")
	if err != nil {
		t.Fatalf("respond survey A: %v", err)
	}
	if satisfactionA.TenantID != "tenant-a" || satisfactionA.WorkspaceID != "workspace-a" {
		t.Fatalf("unexpected scope on satisfaction A: %+v", satisfactionA)
	}

	surveysA, totalSurveysA, err := svc.ListSurveys(ctxA, &SatisfactionSurveyListRequest{Page: 1, PageSize: 20})
	if err != nil {
		t.Fatalf("list surveys A: %v", err)
	}
	if totalSurveysA != 1 || len(surveysA) != 1 || surveysA[0].WorkspaceID != "workspace-a" {
		t.Fatalf("unexpected scoped surveys: total=%d surveys=%+v", totalSurveysA, surveysA)
	}

	satisfactionsA, totalSatisfactionsA, err := svc.ListSatisfactions(ctxA, &SatisfactionListRequest{Page: 1, PageSize: 20})
	if err != nil {
		t.Fatalf("list satisfactions A: %v", err)
	}
	if totalSatisfactionsA != 1 || len(satisfactionsA) != 1 || satisfactionsA[0].WorkspaceID != "workspace-a" {
		t.Fatalf("unexpected scoped satisfactions: total=%d items=%+v", totalSatisfactionsA, satisfactionsA)
	}

	statsA, err := svc.GetSatisfactionStats(ctxA, nil, nil)
	if err != nil {
		t.Fatalf("stats A: %v", err)
	}
	if statsA.TotalRatings != 1 || statsA.RatingDistribution[5] != 1 {
		t.Fatalf("unexpected scoped stats: %+v", statsA)
	}

	if _, err := svc.ResendSurvey(ctxB, surveyA.ID); err == nil {
		t.Fatal("expected cross-workspace resend to fail")
	}
	if _, err := svc.UpdateSatisfaction(ctxB, satisfactionA.ID, "cross-workspace"); err == nil {
		t.Fatal("expected cross-workspace update to fail")
	}
	if err := svc.DeleteSatisfaction(ctxB, satisfactionA.ID); err == nil {
		t.Fatal("expected cross-workspace delete to fail")
	}
}
