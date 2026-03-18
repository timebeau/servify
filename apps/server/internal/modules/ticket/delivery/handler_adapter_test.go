package delivery

import (
	"context"
	"errors"
	"testing"

	ticketapp "servify/apps/server/internal/modules/ticket/application"
)

type fakeTicketStatsQuery struct {
	stats *ticketapp.TicketStatsDTO
	err   error
}

func (f *fakeTicketStatsQuery) GetTicketStats(ctx context.Context, agentID *uint) (*ticketapp.TicketStatsDTO, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.stats, nil
}

func TestHandlerServiceAdapterGetTicketStatsMapsDTO(t *testing.T) {
	a := &HandlerServiceAdapter{
		query: &fakeTicketStatsQuery{
			stats: &ticketapp.TicketStatsDTO{
				Total:        10,
				TodayCreated: 2,
				Pending:      4,
				Resolved:     3,
				ByStatus: []ticketapp.StatusCountDTO{
					{Status: "open", Count: 4},
					{Status: "resolved", Count: 3},
				},
				ByPriority: []ticketapp.PriorityCountDTO{
					{Priority: "high", Count: 5},
				},
			},
		},
	}

	stats, err := a.GetTicketStats(context.Background(), nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if stats.Total != 10 || stats.TodayCreated != 2 || stats.Pending != 4 || stats.Resolved != 3 {
		t.Fatalf("unexpected top-level stats: %+v", stats)
	}
	if len(stats.ByStatus) != 2 || stats.ByStatus[0].Status != "open" || stats.ByStatus[0].Count != 4 {
		t.Fatalf("unexpected status stats: %+v", stats.ByStatus)
	}
	if len(stats.ByPriority) != 1 || stats.ByPriority[0].Priority != "high" || stats.ByPriority[0].Count != 5 {
		t.Fatalf("unexpected priority stats: %+v", stats.ByPriority)
	}
}

func TestHandlerServiceAdapterGetTicketStatsReturnsError(t *testing.T) {
	a := &HandlerServiceAdapter{
		query: &fakeTicketStatsQuery{err: errors.New("query failed")},
	}

	_, err := a.GetTicketStats(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error")
	}
}
