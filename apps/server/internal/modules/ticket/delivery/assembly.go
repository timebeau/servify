package delivery

import (
	"context"

	"servify/apps/server/internal/models"
	ticketapp "servify/apps/server/internal/modules/ticket/application"
	ticketinfra "servify/apps/server/internal/modules/ticket/infra"
	ticketorchestration "servify/apps/server/internal/modules/ticket/orchestration"
	"servify/apps/server/internal/platform/eventbus"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type HandlerAssemblyDependencies struct {
	DB           *gorm.DB
	Logger       *logrus.Logger
	Bus          eventbus.Bus
	SLA          ticketorchestration.SLAService
	Satisfaction ticketorchestration.SatisfactionService
}

func NewHandlerServiceWithDependencies(deps HandlerAssemblyDependencies) *HandlerServiceAdapter {
	repo := ticketinfra.NewGormRepository(deps.DB)
	cmd := ticketapp.NewCommandServiceWithBus(repo, deps.Bus)
	adapter := &HandlerServiceAdapter{
		repo:  repo,
		query: ticketapp.NewQueryService(repo),
		cmd:   cmd,
	}
	orchestrator := ticketorchestration.NewTicketOrchestrator(
		deps.DB,
		deps.Logger,
		deps.SLA,
		deps.Satisfaction,
		deps.Bus,
		repo.CustomerExists,
		repo.FindAutoAssignableAgent,
		func(ctx context.Context, provided map[string]interface{}, ticketCtx map[string]interface{}, enforceRequired bool) ([]models.TicketCustomFieldValue, error) {
			fields, err := repo.ListTicketCustomFields(ctx, true)
			if err != nil {
				return nil, err
			}
			return ticketapp.BuildModelCustomFieldValues(fields, provided, ticketCtx, enforceRequired)
		},
		func(ctx context.Context, ticketID uint, provided map[string]interface{}, ticketCtx map[string]interface{}) (*ticketapp.CustomFieldMutation, error) {
			fields, err := repo.ListTicketCustomFields(ctx, false)
			if err != nil {
				return nil, err
			}
			return ticketapp.PrepareCustomFieldMutation(fields, ticketID, provided, ticketCtx)
		},
		func(ctx context.Context, ticketID uint) (*models.Ticket, error) {
			return repo.LoadTicketModelByID(ctx, ticketID)
		},
		func(ctx context.Context, ticketID uint, agentID uint, assignerID uint) error {
			return adapter.AssignTicket(ctx, ticketID, agentID, assignerID)
		},
		func(ctx context.Context, ticketID uint, userID uint, content string, commentType string) (*models.TicketComment, error) {
			return adapter.AddComment(ctx, ticketID, userID, content, commentType)
		},
	)
	adapter.orchestrator = orchestrator
	return adapter
}
