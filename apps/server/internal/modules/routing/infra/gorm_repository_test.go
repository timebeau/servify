package infra

import (
	"testing"
	"time"

	"servify/apps/server/internal/models"
	"servify/apps/server/internal/modules/routing/domain"
)

func TestMapTransferRecordModel(t *testing.T) {
	now := time.Now()
	fromAgentID := uint(1)
	model := mapTransferRecordModel(domain.Assignment{
		SessionID:   "sess-1",
		FromAgentID: &fromAgentID,
		ToAgentID:   9,
		Reason:      "handoff",
		Notes:       "vip",
		AssignedAt:  now,
	})

	if model.SessionID != "sess-1" || model.ToAgentID == nil || *model.ToAgentID != 9 {
		t.Fatalf("unexpected transfer record model: %+v", model)
	}
	if model.FromAgentID == nil || *model.FromAgentID != 1 || model.Reason != "handoff" {
		t.Fatalf("unexpected from agent/reason mapping: %+v", model)
	}
}

func TestMapWaitingRecordModelAndBack(t *testing.T) {
	now := time.Now()
	assignedAt := now.Add(time.Minute)
	assignedTo := uint(7)

	model := mapWaitingRecordModel(domain.QueueEntry{
		SessionID:    "sess-1",
		Reason:       "no_agent",
		TargetSkills: []string{"billing", "vip"},
		Priority:     "high",
		Notes:        "night shift",
		Status:       domain.QueueStatusTransferred,
		QueuedAt:     now,
		AssignedAt:   &assignedAt,
		AssignedTo:   &assignedTo,
	})
	if model.SessionID != "sess-1" || model.Status != "transferred" || model.TargetSkills == "" {
		t.Fatalf("unexpected waiting record model: %+v", model)
	}

	got := mapQueueEntry(models.WaitingRecord{
		SessionID:    "sess-1",
		Reason:       "no_agent",
		TargetSkills: model.TargetSkills,
		Priority:     "high",
		Notes:        "night shift",
		Status:       "transferred",
		QueuedAt:     now,
		AssignedAt:   &assignedAt,
		AssignedTo:   &assignedTo,
	})
	if got.Status != domain.QueueStatusTransferred || len(got.TargetSkills) != 2 || got.TargetSkills[0] != "billing" {
		t.Fatalf("unexpected queue entry mapping: %+v", got)
	}
	if got.AssignedTo == nil || *got.AssignedTo != 7 {
		t.Fatalf("expected assigned_to mapping, got %+v", got.AssignedTo)
	}
}

func TestUnmarshalSkillsFallbackToCSV(t *testing.T) {
	got := unmarshalSkills("billing, vip,enterprise")
	if len(got) != 3 || got[1] != "vip" {
		t.Fatalf("unexpected skill parsing: %+v", got)
	}
}
