package delivery

import (
	"context"

	"servify/apps/server/internal/models"
	"servify/apps/server/internal/services"
)

// HandlerServiceAdapter bridges the legacy customer service to the handler-facing contract.
type HandlerServiceAdapter struct {
	service *services.CustomerService
}

func NewHandlerServiceAdapter(service *services.CustomerService) *HandlerServiceAdapter {
	return &HandlerServiceAdapter{service: service}
}

func (a *HandlerServiceAdapter) CreateCustomer(ctx context.Context, req *services.CustomerCreateRequest) (*models.User, error) {
	return a.service.CreateCustomer(ctx, req)
}

func (a *HandlerServiceAdapter) GetCustomerByID(ctx context.Context, customerID uint) (*models.User, error) {
	return a.service.GetCustomerByID(ctx, customerID)
}

func (a *HandlerServiceAdapter) UpdateCustomer(ctx context.Context, customerID uint, req *services.CustomerUpdateRequest) (*models.User, error) {
	return a.service.UpdateCustomer(ctx, customerID, req)
}

func (a *HandlerServiceAdapter) ListCustomers(ctx context.Context, req *services.CustomerListRequest) ([]services.CustomerInfo, int64, error) {
	return a.service.ListCustomers(ctx, req)
}

func (a *HandlerServiceAdapter) GetCustomerActivity(ctx context.Context, customerID uint, limit int) (*services.CustomerActivity, error) {
	return a.service.GetCustomerActivity(ctx, customerID, limit)
}

func (a *HandlerServiceAdapter) AddCustomerNote(ctx context.Context, customerID uint, note string, userID uint) error {
	return a.service.AddCustomerNote(ctx, customerID, note, userID)
}

func (a *HandlerServiceAdapter) UpdateCustomerTags(ctx context.Context, customerID uint, tags []string) error {
	return a.service.UpdateCustomerTags(ctx, customerID, tags)
}

func (a *HandlerServiceAdapter) GetCustomerStats(ctx context.Context) (*services.CustomerStats, error) {
	return a.service.GetCustomerStats(ctx)
}
