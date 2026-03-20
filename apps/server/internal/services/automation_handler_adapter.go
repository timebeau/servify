package services

import (
	"context"

	"servify/apps/server/internal/models"
	automationapp "servify/apps/server/internal/modules/automation/application"
	automationdelivery "servify/apps/server/internal/modules/automation/delivery"
)

// AutomationHandlerAdapter bridges the legacy automation service to the handler-facing contract.
type AutomationHandlerAdapter struct {
	service *AutomationService
}

func NewAutomationHandlerAdapter(service *AutomationService) *AutomationHandlerAdapter {
	return &AutomationHandlerAdapter{service: service}
}

var _ automationdelivery.HandlerService = (*AutomationHandlerAdapter)(nil)

func (a *AutomationHandlerAdapter) ListTriggers(ctx context.Context) ([]models.AutomationTrigger, error) {
	return a.service.ListTriggers(ctx)
}

func (a *AutomationHandlerAdapter) CreateTrigger(ctx context.Context, req *automationapp.TriggerRequest) (*models.AutomationTrigger, error) {
	return a.service.CreateTrigger(ctx, req)
}

func (a *AutomationHandlerAdapter) DeleteTrigger(ctx context.Context, id uint) error {
	return a.service.DeleteTrigger(ctx, id)
}

func (a *AutomationHandlerAdapter) ListRuns(ctx context.Context, req *automationapp.RunListQuery) ([]models.AutomationRun, int64, error) {
	return a.service.ListRuns(ctx, req)
}

func (a *AutomationHandlerAdapter) BatchRun(ctx context.Context, req *automationapp.BatchRunRequest) (*automationapp.BatchRunResponse, error) {
	return a.service.BatchRun(ctx, req)
}
