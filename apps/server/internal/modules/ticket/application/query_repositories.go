package application

import (
	"context"

	"servify/apps/server/internal/modules/ticket/domain"
)

type QueryRepository interface {
	GetTicketByID(ctx context.Context, ticketID uint) (*domain.TicketDetails, error)
	ListTickets(ctx context.Context, query ListTicketsQuery) ([]domain.Ticket, int64, error)
	GetTicketStats(ctx context.Context, agentID *uint) (*TicketStatsDTO, error)
}
