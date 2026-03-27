package application

import (
	"context"
	"time"

	"servify/apps/server/internal/models"
)

type Repository interface {
	CreateCustomer(ctx context.Context, cmd CreateCustomerCommand) (*models.User, error)
	GetCustomerByID(ctx context.Context, customerID uint) (*models.User, error)
	UpdateCustomer(ctx context.Context, customerID uint, cmd UpdateCustomerCommand) (*models.User, error)
	ListCustomers(ctx context.Context, query ListCustomersQuery) ([]CustomerInfoDTO, int64, error)
	GetCustomerActivity(ctx context.Context, customerID uint, limit int) (*CustomerActivityDTO, error)
	AddNote(ctx context.Context, customerID uint, note CustomerNoteDTO) error
	UpdateTags(ctx context.Context, customerID uint, tags []string) error
	GetStats(ctx context.Context) (*CustomerStatsDTO, error)
	RevokeCustomerTokens(ctx context.Context, customerID uint, revokeAt time.Time) (int, error)
}
