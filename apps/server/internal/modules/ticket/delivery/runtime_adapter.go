package delivery

import (
	"context"

	ticketapp "servify/apps/server/internal/modules/ticket/application"
	ticketinfra "servify/apps/server/internal/modules/ticket/infra"
	"servify/apps/server/internal/platform/eventbus"

	"gorm.io/gorm"
)

type RuntimeAdapter struct {
	bus eventbus.Bus
}

func NewRuntimeAdapter(bus eventbus.Bus) *RuntimeAdapter {
	return &RuntimeAdapter{bus: bus}
}

func (a *RuntimeAdapter) SyncTransferAssignment(ctx context.Context, tx *gorm.DB, ticketID uint, agentID uint, actorID uint) error {
	repo := ticketinfra.NewGormRepository(tx)
	current, err := repo.GetTicket(ctx, ticketID)
	if err != nil {
		return err
	}

	cmd := ticketapp.UpdateTicketCommand{
		AgentID: &agentID,
		ActorID: actorID,
	}
	if current.Status == "" || current.Status == "open" {
		status := "assigned"
		cmd.Status = &status
	}

	service := ticketapp.NewCommandServiceWithBus(repo, a.bus)
	_, err = service.UpdateTicket(ctx, ticketID, cmd)
	return err
}

var _ RuntimeService = (*RuntimeAdapter)(nil)
