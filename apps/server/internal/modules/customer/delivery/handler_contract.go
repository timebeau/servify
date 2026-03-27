package delivery

import (
	"context"

	"servify/apps/server/internal/models"
	"servify/apps/server/internal/services"
)

// HandlerService is the only customer contract that HTTP handlers should depend on.
type HandlerService interface {
	CreateCustomer(ctx context.Context, req *services.CustomerCreateRequest) (*models.User, error)
	GetCustomerByID(ctx context.Context, customerID uint) (*models.User, error)
	UpdateCustomer(ctx context.Context, customerID uint, req *services.CustomerUpdateRequest) (*models.User, error)
	ListCustomers(ctx context.Context, req *services.CustomerListRequest) ([]services.CustomerInfo, int64, error)
	GetCustomerActivity(ctx context.Context, customerID uint, limit int) (*services.CustomerActivity, error)
	AddCustomerNote(ctx context.Context, customerID uint, note string, userID uint) error
	UpdateCustomerTags(ctx context.Context, customerID uint, tags []string) error
	GetCustomerStats(ctx context.Context) (*services.CustomerStats, error)
	RevokeCustomerTokens(ctx context.Context, customerID uint) (int, error)
}
