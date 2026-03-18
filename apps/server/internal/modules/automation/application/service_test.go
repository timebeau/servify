package application

import (
	"context"
	"testing"

	"servify/apps/server/internal/models"
)

type stubRepo struct {
	triggers []models.AutomationTrigger
	runs     []string
	priority string
}

func (s *stubRepo) ListTriggers(ctx context.Context) ([]models.AutomationTrigger, error) {
	return s.triggers, nil
}
func (s *stubRepo) ListActiveTriggersByEvent(ctx context.Context, event string) ([]models.AutomationTrigger, error) {
	return s.triggers, nil
}
func (s *stubRepo) CreateTrigger(ctx context.Context, req TriggerRequest) (*models.AutomationTrigger, error) {
	return &models.AutomationTrigger{ID: 1, Name: req.Name, Event: req.Event}, nil
}
func (s *stubRepo) DeleteTrigger(ctx context.Context, id uint) error { return nil }
func (s *stubRepo) ListRuns(ctx context.Context, query RunListQuery) ([]models.AutomationRun, int64, error) {
	return nil, 0, nil
}
func (s *stubRepo) RecordRun(ctx context.Context, triggerID uint, ticketID uint, status, message string) error {
	s.runs = append(s.runs, status)
	return nil
}
func (s *stubRepo) GetTicket(ctx context.Context, ticketID uint) (*models.Ticket, error) {
	return &models.Ticket{ID: ticketID, Priority: "normal", Status: "open", Tags: "base"}, nil
}
func (s *stubRepo) UpdateTicketPriority(ctx context.Context, ticketID uint, priority string) error {
	s.priority = priority
	return nil
}
func (s *stubRepo) UpdateTicketTags(ctx context.Context, ticketID uint, tags string) error {
	return nil
}
func (s *stubRepo) CreateTicketComment(ctx context.Context, ticketID uint, content string) error {
	return nil
}

func TestBatchRunDryRunMatches(t *testing.T) {
	repo := &stubRepo{
		triggers: []models.AutomationTrigger{{
			ID:         7,
			Name:       "raise",
			Event:      "ticket.updated",
			Conditions: `[{"field":"ticket.status","op":"eq","value":"open"}]`,
			Actions:    `[{"type":"set_priority","params":{"priority":"high"}}]`,
			Active:     true,
		}},
	}
	svc := NewService(repo)
	resp, err := svc.BatchRun(context.Background(), BatchRunRequest{
		Event:     "ticket.updated",
		TicketIDs: []uint{1},
		DryRun:    true,
	})
	if err != nil {
		t.Fatalf("BatchRun() error = %v", err)
	}
	if resp.Matches != 1 {
		t.Fatalf("expected 1 match, got %d", resp.Matches)
	}
	if repo.priority != "" {
		t.Fatalf("dry run should not update priority")
	}
}
