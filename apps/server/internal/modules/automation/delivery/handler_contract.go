package delivery

import (
	"context"

	"servify/apps/server/internal/models"
	automationapp "servify/apps/server/internal/modules/automation/application"
)

// HandlerService is the only automation contract that HTTP handlers should depend on.
type HandlerService interface {
	ListTriggers(ctx context.Context) ([]models.AutomationTrigger, error)
	CreateTrigger(ctx context.Context, req *automationapp.TriggerRequest) (*models.AutomationTrigger, error)
	DeleteTrigger(ctx context.Context, id uint) error
	ListRuns(ctx context.Context, req *automationapp.RunListQuery) ([]models.AutomationRun, int64, error)
	BatchRun(ctx context.Context, req *automationapp.BatchRunRequest) (*automationapp.BatchRunResponse, error)
}
