package application

import "servify/apps/server/internal/modules/ticket/domain"

type ListTicketsQuery struct {
	Page               int
	PageSize           int
	Status             []string
	Priority           []string
	Category           []string
	AgentID            *uint
	CustomerID         *uint
	Search             string
	SortBy             string
	SortOrder          string
	CustomFieldFilters map[string]string
}

type ListTicketsResult struct {
	Items []domain.Ticket
	Total int64
}
