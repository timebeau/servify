package services

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"servify/apps/server/internal/models"
	conversationdelivery "servify/apps/server/internal/modules/conversation/delivery"
	routingdelivery "servify/apps/server/internal/modules/routing/delivery"

	"gorm.io/gorm"
)

type stubAIForTransferUnit struct{}

func (s stubAIForTransferUnit) ProcessQuery(ctx context.Context, query string, sessionID string) (*AIResponse, error) {
	return &AIResponse{Content: "ok", Confidence: 1, Source: "ai"}, nil
}

func (s stubAIForTransferUnit) ShouldTransferToHuman(query string, sessionHistory []models.Message) bool {
	return false
}

func (s stubAIForTransferUnit) GetSessionSummary(messages []models.Message) (string, error) {
	return "summary", nil
}

func (s stubAIForTransferUnit) InitializeKnowledgeBase() {}

func (s stubAIForTransferUnit) GetStatus(ctx context.Context) map[string]interface{} {
	return map[string]interface{}{"type": "stub"}
}

type stubConversationRuntime struct {
	loadSession func(ctx context.Context, sessionID string) (*conversationdelivery.TransferSession, error)
}

type stubRoutingRuntime struct {
	getTransferHistory   func(ctx context.Context, sessionID string) ([]models.TransferRecord, error)
	listWaitingRecords   func(ctx context.Context, status string, limit int) ([]models.WaitingRecord, error)
	getWaitingRecord     func(ctx context.Context, sessionID string) (*models.WaitingRecord, error)
	cancelWaiting        func(ctx context.Context, tx *gorm.DB, sessionID string, reason string) (*models.WaitingRecord, error)
	markWaitingTransferred func(ctx context.Context, tx *gorm.DB, sessionID string, agentID uint, assignedAt time.Time) (*models.WaitingRecord, error)
}

func (s stubConversationRuntime) LoadTransferSession(ctx context.Context, sessionID string) (*conversationdelivery.TransferSession, error) {
	return s.loadSession(ctx, sessionID)
}

func (s stubConversationRuntime) SyncTransferAssignment(ctx context.Context, tx *gorm.DB, sessionID string, customerID uint, agentID uint) error {
	return nil
}

func (s stubConversationRuntime) SyncWaitingAssignment(ctx context.Context, tx *gorm.DB, sessionID string, customerID uint) error {
	return nil
}

func (s stubConversationRuntime) AppendSystemMessage(ctx context.Context, tx *gorm.DB, sessionID string, content string, createdAt time.Time) error {
	return nil
}

func (s stubRoutingRuntime) AddToWaitingQueue(ctx context.Context, tx *gorm.DB, sessionID string, reason string, targetSkills []string, priority string, notes string) (*models.WaitingRecord, error) {
	return nil, nil
}

func (s stubRoutingRuntime) AssignAgent(ctx context.Context, tx *gorm.DB, cmd routingdelivery.AssignAgentCommand) (*models.TransferRecord, error) {
	return nil, nil
}

func (s stubRoutingRuntime) GetTransferHistory(ctx context.Context, sessionID string) ([]models.TransferRecord, error) {
	if s.getTransferHistory == nil {
		return nil, nil
	}
	return s.getTransferHistory(ctx, sessionID)
}

func (s stubRoutingRuntime) ListWaitingRecords(ctx context.Context, status string, limit int) ([]models.WaitingRecord, error) {
	if s.listWaitingRecords == nil {
		return nil, nil
	}
	return s.listWaitingRecords(ctx, status, limit)
}

func (s stubRoutingRuntime) GetWaitingRecord(ctx context.Context, sessionID string) (*models.WaitingRecord, error) {
	if s.getWaitingRecord == nil {
		return nil, nil
	}
	return s.getWaitingRecord(ctx, sessionID)
}

func (s stubRoutingRuntime) CancelWaiting(ctx context.Context, tx *gorm.DB, sessionID string, reason string) (*models.WaitingRecord, error) {
	if s.cancelWaiting == nil {
		return nil, nil
	}
	return s.cancelWaiting(ctx, tx, sessionID, reason)
}

func (s stubRoutingRuntime) MarkWaitingTransferred(ctx context.Context, tx *gorm.DB, sessionID string, agentID uint, assignedAt time.Time) (*models.WaitingRecord, error) {
	if s.markWaitingTransferred == nil {
		return nil, nil
	}
	return s.markWaitingTransferred(ctx, tx, sessionID, agentID, assignedAt)
}

func TestSessionTransferService_LoadTransferSession_UsesConversationAdapterError(t *testing.T) {
	expectedErr := errors.New("conversation adapter failed")
	svc := NewSessionTransferServiceWithAdapters(nil, nil, stubAIForTransferUnit{}, nil, nil, SessionTransferAdapters{
		Conversation: stubConversationRuntime{
			loadSession: func(ctx context.Context, sessionID string) (*conversationdelivery.TransferSession, error) {
				if sessionID != "s-adapter" {
					t.Fatalf("unexpected session id: %s", sessionID)
				}
				return nil, expectedErr
			},
		},
	})

	_, err := svc.loadTransferSession(context.Background(), "s-adapter")
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected adapter error, got %v", err)
	}
}

func TestNormalizeWaitingRecordQuery(t *testing.T) {
	tests := []struct {
		name       string
		status     string
		limit      int
		wantStatus string
		wantLimit  int
	}{
		{name: "defaults", status: "", limit: 0, wantStatus: "waiting", wantLimit: 50},
		{name: "caps large limit", status: "transferred", limit: 500, wantStatus: "transferred", wantLimit: 50},
		{name: "keeps valid values", status: "cancelled", limit: 25, wantStatus: "cancelled", wantLimit: 25},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStatus, gotLimit := normalizeWaitingRecordQuery(tt.status, tt.limit)
			if gotStatus != tt.wantStatus || gotLimit != tt.wantLimit {
				t.Fatalf("normalizeWaitingRecordQuery(%q, %d) = (%q, %d), want (%q, %d)",
					tt.status, tt.limit, gotStatus, gotLimit, tt.wantStatus, tt.wantLimit)
			}
		})
	}
}

func TestSessionTransferService_ListTransferHistory_WrapsRoutingError(t *testing.T) {
	expectedErr := errors.New("routing unavailable")
	svc := NewSessionTransferServiceWithAdapters(nil, nil, stubAIForTransferUnit{}, nil, nil, SessionTransferAdapters{
		Routing: stubRoutingRuntime{
			getTransferHistory: func(ctx context.Context, sessionID string) ([]models.TransferRecord, error) {
				if sessionID != "s-history" {
					t.Fatalf("unexpected session id: %s", sessionID)
				}
				return nil, expectedErr
			},
		},
	})

	_, err := svc.listTransferHistory(context.Background(), "s-history")
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected wrapped routing error, got %v", err)
	}
	if !strings.Contains(err.Error(), "failed to get transfer history") {
		t.Fatalf("expected wrapped message, got %v", err)
	}
}

func TestSessionTransferService_GetActiveWaitingRecord_NonWaitingIsNotFound(t *testing.T) {
	svc := NewSessionTransferServiceWithAdapters(nil, nil, stubAIForTransferUnit{}, nil, nil, SessionTransferAdapters{
		Routing: stubRoutingRuntime{
			getWaitingRecord: func(ctx context.Context, sessionID string) (*models.WaitingRecord, error) {
				if sessionID != "s-waiting" {
					t.Fatalf("unexpected session id: %s", sessionID)
				}
				return &models.WaitingRecord{SessionID: sessionID, Status: "cancelled"}, nil
			},
		},
	})

	_, err := svc.getActiveWaitingRecord(context.Background(), "s-waiting")
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected ErrRecordNotFound, got %v", err)
	}
}

func TestSessionTransferService_CancelWaitingRecord_WrapsRoutingError(t *testing.T) {
	expectedErr := errors.New("cancel failed")
	svc := NewSessionTransferServiceWithAdapters(nil, nil, stubAIForTransferUnit{}, nil, nil, SessionTransferAdapters{
		Routing: stubRoutingRuntime{
			cancelWaiting: func(ctx context.Context, tx *gorm.DB, sessionID string, reason string) (*models.WaitingRecord, error) {
				if sessionID != "s-cancel" || reason != "user_left" {
					t.Fatalf("unexpected cancel inputs: %s %s", sessionID, reason)
				}
				return nil, expectedErr
			},
		},
	})

	err := svc.cancelWaitingRecord(context.Background(), nil, "s-cancel", "user_left")
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected wrapped error, got %v", err)
	}
	if !strings.Contains(err.Error(), "update waiting record") {
		t.Fatalf("expected wrapped message, got %v", err)
	}
}

func TestBuildTransferTicketUpdate_AssignsOpenTicket(t *testing.T) {
	updates, fromStatus, toStatus := buildTransferTicketUpdate(7, "open")
	if fromStatus != "open" || toStatus != "assigned" {
		t.Fatalf("unexpected statuses: from=%s to=%s", fromStatus, toStatus)
	}
	if updates["agent_id"] != uint(7) {
		t.Fatalf("expected agent_id update, got %+v", updates)
	}
	if updates["status"] != "assigned" {
		t.Fatalf("expected assigned status update, got %+v", updates)
	}
}

func TestBuildTransferTicketUpdate_PreservesNonOpenStatus(t *testing.T) {
	updates, fromStatus, toStatus := buildTransferTicketUpdate(9, "pending")
	if fromStatus != "pending" || toStatus != "pending" {
		t.Fatalf("unexpected statuses: from=%s to=%s", fromStatus, toStatus)
	}
	if updates["agent_id"] != uint(9) {
		t.Fatalf("expected agent_id update, got %+v", updates)
	}
	if _, ok := updates["status"]; ok {
		t.Fatalf("did not expect status update, got %+v", updates)
	}
}

func TestMarkWaitingTransferred_IgnoresRoutingNotFound(t *testing.T) {
	svc := NewSessionTransferServiceWithAdapters(nil, nil, stubAIForTransferUnit{}, nil, nil, SessionTransferAdapters{
		Routing: stubRoutingRuntime{
			markWaitingTransferred: func(ctx context.Context, tx *gorm.DB, sessionID string, agentID uint, assignedAt time.Time) (*models.WaitingRecord, error) {
				if sessionID != "s-transfer" || agentID != 11 {
					t.Fatalf("unexpected transfer inputs: %s %d", sessionID, agentID)
				}
				return nil, gorm.ErrRecordNotFound
			},
		},
	})

	if err := svc.markWaitingTransferred(context.Background(), nil, "s-transfer", 11, time.Now()); err != nil {
		t.Fatalf("expected not found to be ignored, got %v", err)
	}
}
