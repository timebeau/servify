package application

import (
	"context"
	"fmt"
	"testing"
	"time"

	"servify/apps/server/internal/modules/ticket/domain"
)

type stubQueryRepo struct {
	details *domain.TicketDetails
	items   []domain.Ticket
	total   int64
	err     error
}

func (s stubQueryRepo) GetTicketByID(ctx context.Context, ticketID uint) (*domain.TicketDetails, error) {
	if s.err != nil {
		return nil, s.err
	}
	if s.details == nil {
		return nil, fmt.Errorf("not found")
	}
	return s.details, nil
}

func (s stubQueryRepo) ListTickets(ctx context.Context, query ListTicketsQuery) ([]domain.Ticket, int64, error) {
	if s.err != nil {
		return nil, 0, s.err
	}
	return s.items, s.total, nil
}

func TestQueryServiceGetTicketByID(t *testing.T) {
	now := time.Now()
	svc := NewQueryService(stubQueryRepo{
		details: &domain.TicketDetails{
			Ticket: domain.Ticket{
				ID:         1,
				Title:      "Billing issue",
				Description:"Need help",
				CustomerID: 10,
				Status:     "open",
				Priority:   "high",
				Category:   "billing",
				Source:     "web",
				CreatedAt:  now,
				UpdatedAt:  now,
			},
		},
	})

	got, err := svc.GetTicketByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.ID != 1 || got.Title != "Billing issue" {
		t.Fatalf("unexpected ticket details: %+v", got)
	}
}

func TestQueryServiceListTickets(t *testing.T) {
	now := time.Now()
	svc := NewQueryService(stubQueryRepo{
		items: []domain.Ticket{
			{
				ID:         1,
				Title:      "Billing issue",
				CustomerID: 10,
				Status:     "open",
				Priority:   "high",
				Category:   "billing",
				Source:     "web",
				CreatedAt:  now,
				UpdatedAt:  now,
			},
		},
		total: 1,
	})

	got, err := svc.ListTickets(context.Background(), ListTicketsQuery{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.Total != 1 || len(got.Items) != 1 {
		t.Fatalf("unexpected list result: %+v", got)
	}
}
