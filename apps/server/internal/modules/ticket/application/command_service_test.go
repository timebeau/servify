package application

import (
	"context"
	"fmt"
	"testing"
	"time"

	"servify/apps/server/internal/modules/ticket/domain"
)

type stubCommandRepo struct {
	tickets        map[uint]*domain.Ticket
	comments       map[uint][]domain.Comment
	statusChanges  map[uint][]domain.StatusChange
	nextID         uint
	nextCommentID  uint
	nextChangeID   uint
	customerExists bool
	agentAvailable bool
	err            error
}

func (s *stubCommandRepo) CreateTicket(ctx context.Context, ticket *domain.Ticket) error {
	if s.err != nil {
		return s.err
	}
	if s.tickets == nil {
		s.tickets = map[uint]*domain.Ticket{}
	}
	if s.nextID == 0 {
		s.nextID = 1
	}
	ticket.ID = s.nextID
	s.nextID++
	cp := *ticket
	s.tickets[ticket.ID] = &cp
	return nil
}

func (s *stubCommandRepo) UpdateTicket(ctx context.Context, ticket *domain.Ticket) error {
	if s.err != nil {
		return s.err
	}
	cp := *ticket
	s.tickets[ticket.ID] = &cp
	return nil
}

func (s *stubCommandRepo) AddComment(ctx context.Context, ticketID uint, comment *domain.Comment) error {
	if s.err != nil {
		return s.err
	}
	if s.comments == nil {
		s.comments = map[uint][]domain.Comment{}
	}
	if s.nextCommentID == 0 {
		s.nextCommentID = 1
	}
	comment.ID = s.nextCommentID
	s.nextCommentID++
	s.comments[ticketID] = append(s.comments[ticketID], *comment)
	return nil
}

func (s *stubCommandRepo) RecordStatusChange(ctx context.Context, ticketID uint, change *domain.StatusChange) error {
	if s.err != nil {
		return s.err
	}
	if s.statusChanges == nil {
		s.statusChanges = map[uint][]domain.StatusChange{}
	}
	if s.nextChangeID == 0 {
		s.nextChangeID = 1
	}
	change.ID = s.nextChangeID
	s.nextChangeID++
	s.statusChanges[ticketID] = append(s.statusChanges[ticketID], *change)
	return nil
}

func (s *stubCommandRepo) GetTicket(ctx context.Context, ticketID uint) (*domain.Ticket, error) {
	if s.err != nil {
		return nil, s.err
	}
	ticket, ok := s.tickets[ticketID]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	cp := *ticket
	return &cp, nil
}

func (s *stubCommandRepo) CustomerExists(ctx context.Context, customerID uint) (bool, error) {
	if s.err != nil {
		return false, s.err
	}
	return s.customerExists, nil
}

func (s *stubCommandRepo) AgentAssignable(ctx context.Context, agentID uint) (bool, error) {
	if s.err != nil {
		return false, s.err
	}
	return s.agentAvailable, nil
}

func TestCommandServiceCreateTicket(t *testing.T) {
	svc := NewCommandService(&stubCommandRepo{customerExists: true})
	got, err := svc.CreateTicket(context.Background(), CreateTicketCommand{
		Title:      "Need Help",
		CustomerID: 7,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.ID == 0 || got.Status != "open" || got.Source != "web" {
		t.Fatalf("unexpected created ticket: %+v", got)
	}
}

func TestCommandServiceUpdateTicket(t *testing.T) {
	repo := &stubCommandRepo{
		customerExists: true,
		tickets: map[uint]*domain.Ticket{
			1: {
				ID:         1,
				Title:      "Billing",
				CustomerID: 9,
				Status:     "open",
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
		},
	}
	svc := NewCommandService(repo)
	newStatus := "resolved"
	newTitle := "Billing Updated"

	got, err := svc.UpdateTicket(context.Background(), 1, UpdateTicketCommand{
		Title:   &newTitle,
		Status:  &newStatus,
		ActorID: 99,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.Title != "Billing Updated" || got.Status != "resolved" || got.ResolvedAt == nil {
		t.Fatalf("unexpected updated ticket: %+v", got)
	}
	if len(repo.statusChanges[1]) != 1 || repo.statusChanges[1][0].ToStatus != "resolved" {
		t.Fatalf("expected status history to be recorded, got %+v", repo.statusChanges[1])
	}
}

func TestCommandServiceAssignTicket(t *testing.T) {
	repo := &stubCommandRepo{
		agentAvailable: true,
		tickets: map[uint]*domain.Ticket{
			1: {
				ID:         1,
				Title:      "Billing",
				CustomerID: 9,
				Status:     "open",
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
		},
	}
	svc := NewCommandService(repo)
	got, err := svc.AssignTicket(context.Background(), 1, AssignTicketCommand{AgentID: 88})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.AgentID == nil || *got.AgentID != 88 || got.Status != "assigned" {
		t.Fatalf("unexpected assigned ticket: %+v", got)
	}
	if len(repo.statusChanges[1]) != 1 || repo.statusChanges[1][0].ToStatus != "assigned" {
		t.Fatalf("expected assignment status history, got %+v", repo.statusChanges[1])
	}
}

func TestCommandServiceAddComment(t *testing.T) {
	repo := &stubCommandRepo{
		tickets: map[uint]*domain.Ticket{
			1: {
				ID:         1,
				Title:      "Billing",
				CustomerID: 9,
				Status:     "open",
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
		},
	}
	svc := NewCommandService(repo)
	got, err := svc.AddComment(context.Background(), 1, AddCommentCommand{
		UserID:  7,
		Content: " Need more logs ",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.ID == 0 || got.Content != "Need more logs" || got.Type != "comment" {
		t.Fatalf("unexpected comment: %+v", got)
	}
	if len(repo.comments[1]) != 1 {
		t.Fatalf("expected comment to be persisted, got %+v", repo.comments[1])
	}
}

func TestCommandServiceCloseTicket(t *testing.T) {
	repo := &stubCommandRepo{
		tickets: map[uint]*domain.Ticket{
			1: {
				ID:         1,
				Title:      "Billing",
				CustomerID: 9,
				Status:     "open",
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
		},
	}
	svc := NewCommandService(repo)
	got, err := svc.CloseTicket(context.Background(), 1, CloseTicketCommand{UserID: 8, Reason: "done"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.Status != "closed" || got.ClosedAt == nil {
		t.Fatalf("unexpected closed ticket: %+v", got)
	}
	if len(repo.statusChanges[1]) != 1 || repo.statusChanges[1][0].ToStatus != "closed" {
		t.Fatalf("expected close status history, got %+v", repo.statusChanges[1])
	}
}

func TestStatusTransitionPolicyRejectsInvalidTransition(t *testing.T) {
	repo := &stubCommandRepo{
		tickets: map[uint]*domain.Ticket{
			1: {
				ID:         1,
				Title:      "Billing",
				CustomerID: 9,
				Status:     "resolved",
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
		},
	}
	svc := NewCommandService(repo)
	nextStatus := "assigned"

	_, err := svc.UpdateTicket(context.Background(), 1, UpdateTicketCommand{
		Status:  &nextStatus,
		ActorID: 42,
	})
	if err == nil {
		t.Fatal("expected invalid transition error")
	}
}
