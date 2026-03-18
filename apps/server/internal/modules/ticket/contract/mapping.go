package contract

import (
	"strings"

	ticketapp "servify/apps/server/internal/modules/ticket/application"
)

func NormalizeTicketSortBy(sortBy string) string {
	s := strings.ToLower(strings.TrimSpace(sortBy))
	s = strings.TrimPrefix(s, "tickets.")
	switch s {
	case "id", "created_at", "updated_at", "priority", "status", "category", "due_date", "resolved_at", "closed_at":
		return "tickets." + s
	default:
		return "tickets.created_at"
	}
}

func MapTicketStats(stats *ticketapp.TicketStatsDTO) *TicketStats {
	if stats == nil {
		return nil
	}
	return &TicketStats{
		Total:        stats.Total,
		TodayCreated: stats.TodayCreated,
		Pending:      stats.Pending,
		Resolved:     stats.Resolved,
		ByStatus:     mapStatusCounts(stats.ByStatus),
		ByPriority:   mapPriorityCounts(stats.ByPriority),
	}
}

func mapStatusCounts(items []ticketapp.StatusCountDTO) []StatusCount {
	out := make([]StatusCount, 0, len(items))
	for _, item := range items {
		out = append(out, StatusCount{
			Status: item.Status,
			Count:  item.Count,
		})
	}
	return out
}

func mapPriorityCounts(items []ticketapp.PriorityCountDTO) []PriorityCount {
	out := make([]PriorityCount, 0, len(items))
	for _, item := range items {
		out = append(out, PriorityCount{
			Priority: item.Priority,
			Count:    item.Count,
		})
	}
	return out
}
