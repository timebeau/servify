package application

import (
	"context"

	"servify/apps/server/internal/modules/ticket/domain"
)

type CommandRepository interface {
	CreateTicket(ctx context.Context, ticket *domain.Ticket) error
	UpdateTicket(ctx context.Context, ticket *domain.Ticket) error
	UpdateTicketWithStatus(ctx context.Context, ticket *domain.Ticket, fromStatus string, userID uint, reason string) error
	AssignTicket(ctx context.Context, ticket *domain.Ticket, previousAgentID *uint, fromStatus string, userID uint, reason string) error
	UnassignTicket(ctx context.Context, ticket *domain.Ticket, previousAgentID uint, fromStatus string, userID uint, reason string) error
	CloseTicket(ctx context.Context, ticket *domain.Ticket, fromStatus string, userID uint, reason string) error
	AddComment(ctx context.Context, ticketID uint, comment *domain.Comment) error
	RecordStatusChange(ctx context.Context, ticketID uint, change *domain.StatusChange) error
	GetTicket(ctx context.Context, ticketID uint) (*domain.Ticket, error)
	CustomerExists(ctx context.Context, customerID uint) (bool, error)
	AgentAssignable(ctx context.Context, agentID uint) (bool, error)
}
