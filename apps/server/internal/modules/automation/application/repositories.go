package application

import (
	"context"

	"servify/apps/server/internal/models"
)

type Repository interface {
	ListTriggers(ctx context.Context) ([]models.AutomationTrigger, error)
	ListActiveTriggersByEvent(ctx context.Context, event string) ([]models.AutomationTrigger, error)
	CreateTrigger(ctx context.Context, req TriggerRequest) (*models.AutomationTrigger, error)
	DeleteTrigger(ctx context.Context, id uint) error
	ListRuns(ctx context.Context, query RunListQuery) ([]models.AutomationRun, int64, error)
	RecordRun(ctx context.Context, triggerID uint, ticketID uint, status, message string) error
	GetTicket(ctx context.Context, ticketID uint) (*models.Ticket, error)
	UpdateTicketPriority(ctx context.Context, ticketID uint, priority string) error
	UpdateTicketTags(ctx context.Context, ticketID uint, tags string) error
	CreateTicketComment(ctx context.Context, ticketID uint, content string) error
}
