package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"servify/apps/server/internal/models"
	conversationdelivery "servify/apps/server/internal/modules/conversation/delivery"

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
		name         string
		status       string
		limit        int
		wantStatus   string
		wantLimit    int
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
