package infra

import (
	"strings"
	"testing"
	"time"

	"servify/apps/server/internal/models"
	"servify/apps/server/internal/modules/conversation/domain"
)

func TestMapConversationModel(t *testing.T) {
	now := time.Now()
	customerID := uint(7)
	agentID := uint(9)
	lastMessageAt := now.Add(time.Minute)

	model := mapConversationModel(domain.Conversation{
		ID:         "conv-1",
		CustomerID: &customerID,
		Status:     domain.ConversationStatusTransferred,
		Channel: domain.ChannelBinding{
			Channel: "web",
		},
		Participants: []domain.Participant{
			{ID: "customer:7", UserID: &customerID, Role: domain.ParticipantRoleCustomer},
			{ID: "agent:9", UserID: &agentID, Role: domain.ParticipantRoleAgent},
		},
		StartedAt:     now,
		LastMessageAt: &lastMessageAt,
	})

	if model.ID != "conv-1" || model.UserID != customerID || model.Status != "transferred" || model.Platform != "web" {
		t.Fatalf("unexpected session model: %+v", model)
	}
	if model.AgentID == nil || *model.AgentID != agentID {
		t.Fatalf("expected agent id to be mapped, got %+v", model.AgentID)
	}
	if !model.UpdatedAt.Equal(lastMessageAt) {
		t.Fatalf("expected updated_at to follow last_message_at, got %v", model.UpdatedAt)
	}
}

func TestMapConversation(t *testing.T) {
	now := time.Now()
	agentID := uint(9)
	endedAt := now.Add(2 * time.Minute)
	model := models.Session{
		ID:        "conv-1",
		UserID:    7,
		AgentID:   &agentID,
		Status:    "transferred",
		Platform:  "web",
		StartedAt: now,
		UpdatedAt: now.Add(time.Minute),
		EndedAt:   &endedAt,
		User:      models.User{ID: 7, Name: "Customer"},
		Agent:     &models.User{ID: agentID, Name: "Agent"},
	}

	got := mapConversation(model)
	if got.ID != "conv-1" || got.Status != domain.ConversationStatusTransferred {
		t.Fatalf("unexpected conversation mapping: %+v", got)
	}
	if got.CustomerID == nil || *got.CustomerID != 7 {
		t.Fatalf("expected customer id, got %+v", got.CustomerID)
	}
	if len(got.Participants) != 2 {
		t.Fatalf("expected participants to be derived, got %+v", got.Participants)
	}
	if got.LastMessageAt == nil || !got.LastMessageAt.Equal(model.UpdatedAt) {
		t.Fatalf("expected last_message_at to map from updated_at, got %+v", got.LastMessageAt)
	}
	if got.EndedAt == nil || !got.EndedAt.Equal(endedAt) {
		t.Fatalf("expected ended_at to map, got %+v", got.EndedAt)
	}
}

func TestMapMessageModelAndBack(t *testing.T) {
	now := time.Now()
	model := mapMessageModel(domain.ConversationMessage{
		ConversationID: "conv-1",
		Sender:         domain.ParticipantRoleSystem,
		Kind:           domain.MessageKindSystem,
		Content:        "agent assigned",
		CreatedAt:      now,
	})
	if model.SessionID != "conv-1" || model.Sender != "system" || model.Type != "system" {
		t.Fatalf("unexpected message model: %+v", model)
	}

	got := mapMessage(models.Message{
		ID:        12,
		SessionID: "conv-1",
		Sender:    "agent",
		Type:      "text",
		Content:   "hello",
		CreatedAt: now,
	})
	if got.ID != "12" || got.Sender != domain.ParticipantRoleAgent || got.Kind != domain.MessageKindText || got.Content != "hello" {
		t.Fatalf("unexpected message mapping: %+v", got)
	}
}

func TestMapConversationStatus_Bidirectional(t *testing.T) {
	cases := []struct {
		dbStatus   string
		domainStatus domain.ConversationStatus
	}{
		{"active", domain.ConversationStatusActive},
		{"ended", domain.ConversationStatusClosed},
		{"transferred", domain.ConversationStatusTransferred},
		{"waiting_human", domain.ConversationStatusWaitingHuman},
		{"ACTIVE", domain.ConversationStatusActive},
		{" Waiting_Human ", domain.ConversationStatusWaitingHuman},
	}

	for _, tc := range cases {
		got := mapSessionStatusToConversationStatus(tc.dbStatus)
		if got != tc.domainStatus {
			t.Errorf("mapSessionStatusToConversationStatus(%q) = %q, want %q", tc.dbStatus, got, tc.domainStatus)
		}

		// 反向映射（仅小写精确匹配）
		if tc.dbStatus == strings.TrimSpace(strings.ToLower(tc.dbStatus)) {
			reverse := mapConversationStatusToSessionStatus(tc.domainStatus)
			if reverse != tc.dbStatus {
				t.Errorf("mapConversationStatusToSessionStatus(%q) = %q, want %q", tc.domainStatus, reverse, tc.dbStatus)
			}
		}
	}
}
