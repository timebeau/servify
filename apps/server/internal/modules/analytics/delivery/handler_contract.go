package delivery

import (
	"context"
	"time"

	analyticscontract "servify/apps/server/internal/modules/analytics/contract"
)

// HandlerService is the only analytics contract that HTTP handlers should depend on.
type HandlerService interface {
	GetDashboardStats(ctx context.Context) (*analyticscontract.DashboardStats, error)
	GetTimeRangeStats(ctx context.Context, startDate, endDate time.Time) ([]analyticscontract.TimeRangeStats, error)
	GetAgentPerformanceStats(ctx context.Context, startDate, endDate time.Time, limit int) ([]analyticscontract.AgentPerformanceStats, error)
	GetTicketCategoryStats(ctx context.Context, startDate, endDate time.Time) ([]analyticscontract.CategoryStats, error)
	GetTicketPriorityStats(ctx context.Context, startDate, endDate time.Time) ([]analyticscontract.CategoryStats, error)
	GetCustomerSourceStats(ctx context.Context) ([]analyticscontract.CategoryStats, error)
	GetRemoteAssistTicketStats(ctx context.Context) (*analyticscontract.RemoteAssistTicketStats, error)
	UpdateDailyStats(ctx context.Context, date time.Time) error
}
