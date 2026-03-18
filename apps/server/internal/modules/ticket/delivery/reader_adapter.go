package delivery

import (
	"context"

	"servify/apps/server/internal/models"
	ticketinfra "servify/apps/server/internal/modules/ticket/infra"

	"gorm.io/gorm"
)

// ReaderServiceAdapter exposes only ticket read operations for cross-module consumers.
type ReaderServiceAdapter struct {
	repo *ticketinfra.GormRepository
}

func NewReaderServiceAdapter(db *gorm.DB) *ReaderServiceAdapter {
	return &ReaderServiceAdapter{
		repo: ticketinfra.NewGormRepository(db),
	}
}

func (a *ReaderServiceAdapter) GetTicketByID(ctx context.Context, ticketID uint) (*models.Ticket, error) {
	return a.repo.LoadTicketModelByID(ctx, ticketID)
}
