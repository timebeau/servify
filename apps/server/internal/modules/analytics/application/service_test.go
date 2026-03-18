package application

import (
	"context"
	"testing"
	"time"
)

type stubRepo struct {
	last IncrementEvent
}

func (s *stubRepo) GetDashboardStats(ctx context.Context) (*DashboardStats, error) {
	return &DashboardStats{TotalTickets: 1}, nil
}
func (s *stubRepo) GetTimeRangeStats(ctx context.Context, startDate, endDate time.Time) ([]TimeRangeStats, error) {
	return []TimeRangeStats{{Date: startDate.Format("2006-01-02")}}, nil
}
func (s *stubRepo) GetAgentPerformanceStats(ctx context.Context, startDate, endDate time.Time, limit int) ([]AgentPerformanceStats, error) {
	return []AgentPerformanceStats{{AgentID: 1}}, nil
}
func (s *stubRepo) GetTicketCategoryStats(ctx context.Context, startDate, endDate time.Time) ([]CategoryStats, error) {
	return []CategoryStats{{Category: "general", Count: 1}}, nil
}
func (s *stubRepo) GetTicketPriorityStats(ctx context.Context, startDate, endDate time.Time) ([]CategoryStats, error) {
	return []CategoryStats{{Category: "normal", Count: 1}}, nil
}
func (s *stubRepo) GetCustomerSourceStats(ctx context.Context) ([]CategoryStats, error) {
	return []CategoryStats{{Category: "web", Count: 1}}, nil
}
func (s *stubRepo) UpdateDailyStats(ctx context.Context, date time.Time) error { return nil }
func (s *stubRepo) IncrementDailyStat(ctx context.Context, event IncrementEvent) error {
	s.last = event
	return nil
}

func TestIncrementDailyStatDefaultsDate(t *testing.T) {
	repo := &stubRepo{}
	svc := NewService(repo)
	if err := svc.IncrementDailyStat(context.Background(), IncrementEvent{Kind: IncrementTickets}); err != nil {
		t.Fatalf("IncrementDailyStat() error = %v", err)
	}
	if repo.last.Date.IsZero() {
		t.Fatal("expected date to be defaulted")
	}
	if repo.last.Kind != IncrementTickets {
		t.Fatalf("unexpected kind: %s", repo.last.Kind)
	}
}
