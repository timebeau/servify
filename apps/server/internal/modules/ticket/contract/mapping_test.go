package contract

import (
	"testing"

	ticketapp "servify/apps/server/internal/modules/ticket/application"
)

func TestNormalizeTicketSortBy(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "empty defaults", input: "", want: "tickets.created_at"},
		{name: "plain field", input: "priority", want: "tickets.priority"},
		{name: "qualified field", input: "tickets.updated_at", want: "tickets.updated_at"},
		{name: "trim and lower", input: "  STATUS  ", want: "tickets.status"},
		{name: "unsupported defaults", input: "deleted_at", want: "tickets.created_at"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NormalizeTicketSortBy(tt.input); got != tt.want {
				t.Fatalf("NormalizeTicketSortBy(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestMapTicketStats(t *testing.T) {
	stats := MapTicketStats(&ticketapp.TicketStatsDTO{
		Total:        10,
		TodayCreated: 2,
		Pending:      4,
		Resolved:     3,
		ByStatus: []ticketapp.StatusCountDTO{
			{Status: "open", Count: 4},
		},
		ByPriority: []ticketapp.PriorityCountDTO{
			{Priority: "high", Count: 5},
		},
	})

	if stats == nil {
		t.Fatal("expected mapped stats")
	}
	if stats.Total != 10 || stats.TodayCreated != 2 || stats.Pending != 4 || stats.Resolved != 3 {
		t.Fatalf("unexpected top-level stats: %+v", stats)
	}
	if len(stats.ByStatus) != 1 || stats.ByStatus[0].Status != "open" || stats.ByStatus[0].Count != 4 {
		t.Fatalf("unexpected status mapping: %+v", stats.ByStatus)
	}
	if len(stats.ByPriority) != 1 || stats.ByPriority[0].Priority != "high" || stats.ByPriority[0].Count != 5 {
		t.Fatalf("unexpected priority mapping: %+v", stats.ByPriority)
	}
}
