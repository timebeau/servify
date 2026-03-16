package application

import (
	"context"

	"servify/apps/server/internal/modules/ticket/domain"
)

type CommandRepository interface {
	CreateTicket(ctx context.Context, ticket *domain.Ticket) error
	UpdateTicket(ctx context.Context, ticket *domain.Ticket) error
	AddComment(ctx context.Context, ticketID uint, comment *domain.Comment) error
	RecordStatusChange(ctx context.Context, ticketID uint, change *domain.StatusChange) error
	GetTicket(ctx context.Context, ticketID uint) (*domain.Ticket, error)
	CustomerExists(ctx context.Context, customerID uint) (bool, error)
	AgentAssignable(ctx context.Context, agentID uint) (bool, error)
}
