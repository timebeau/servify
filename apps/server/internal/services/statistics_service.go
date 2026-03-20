package services

import (
	"context"
	"time"

	analyticsapp "servify/apps/server/internal/modules/analytics/application"
	analyticscontract "servify/apps/server/internal/modules/analytics/contract"
	analyticsdelivery "servify/apps/server/internal/modules/analytics/delivery"
	analyticsinfra "servify/apps/server/internal/modules/analytics/infra"
	"servify/apps/server/internal/platform/eventbus"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// StatisticsService 数据统计服务兼容层。
//
// 它目前承担两类职责：
// 1. 作为 HTTP handler 的兼容 contract 实现
// 2. 维护 event bus 订阅与 DTO 映射等尚未完全迁出的 glue logic
//
// 后续迁移中，HTTP 层应只依赖 modules/analytics/delivery.HandlerService。
type StatisticsService struct {
	db         *gorm.DB
	logger     *logrus.Logger
	module     *analyticsapp.Service
	subscriber *analyticsdelivery.EventBusSubscriber
}

// NewStatisticsService 创建统计服务。
func NewStatisticsService(db *gorm.DB, logger *logrus.Logger) *StatisticsService {
	if logger == nil {
		logger = logrus.New()
	}
	repo := analyticsinfra.NewGormRepository(db)
	module := analyticsapp.NewService(repo)
	return &StatisticsService{
		db:         db,
		logger:     logger,
		module:     module,
		subscriber: analyticsdelivery.NewEventBusSubscriber(module),
	}
}

func (s *StatisticsService) SetEventBus(bus eventbus.Bus) {
	if s.subscriber != nil {
		s.subscriber.Register(bus)
	}
}

type DashboardStats = analyticscontract.DashboardStats
type TimeRangeStats = analyticscontract.TimeRangeStats
type AgentPerformanceStats = analyticscontract.AgentPerformanceStats
type CategoryStats = analyticscontract.CategoryStats

// GetDashboardStats 获取仪表板统计数据
func (s *StatisticsService) GetDashboardStats(ctx context.Context) (*DashboardStats, error) {
	stats, err := s.module.GetDashboardStats(ctx)
	if err != nil {
		return nil, err
	}
	return dashboardStatsFromDTO(stats), nil
}

// GetTimeRangeStats 获取时间范围统计
func (s *StatisticsService) GetTimeRangeStats(ctx context.Context, startDate, endDate time.Time) ([]TimeRangeStats, error) {
	stats, err := s.module.GetTimeRangeStats(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}
	return timeRangeStatsFromDTO(stats), nil
}

// GetAgentPerformanceStats 获取客服绩效统计
func (s *StatisticsService) GetAgentPerformanceStats(ctx context.Context, startDate, endDate time.Time, limit int) ([]AgentPerformanceStats, error) {
	stats, err := s.module.GetAgentPerformanceStats(ctx, startDate, endDate, limit)
	if err != nil {
		return nil, err
	}
	return agentPerformanceStatsFromDTO(stats), nil
}

// GetTicketCategoryStats 获取工单分类统计
func (s *StatisticsService) GetTicketCategoryStats(ctx context.Context, startDate, endDate time.Time) ([]CategoryStats, error) {
	stats, err := s.module.GetTicketCategoryStats(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}
	return categoryStatsFromDTO(stats), nil
}

// GetTicketPriorityStats 获取工单优先级统计
func (s *StatisticsService) GetTicketPriorityStats(ctx context.Context, startDate, endDate time.Time) ([]CategoryStats, error) {
	stats, err := s.module.GetTicketPriorityStats(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}
	return categoryStatsFromDTO(stats), nil
}

// GetCustomerSourceStats 获取客户来源统计
func (s *StatisticsService) GetCustomerSourceStats(ctx context.Context) ([]CategoryStats, error) {
	stats, err := s.module.GetCustomerSourceStats(ctx)
	if err != nil {
		return nil, err
	}
	return categoryStatsFromDTO(stats), nil
}

// UpdateDailyStats 更新每日统计数据
func (s *StatisticsService) UpdateDailyStats(ctx context.Context, date time.Time) error {
	return s.module.UpdateDailyStats(ctx, date)
}

// IncrementAIUsage 增加 AI 使用计数
func (s *StatisticsService) IncrementAIUsage(ctx context.Context) {
	_ = s.module.IncrementDailyStat(ctx, analyticsapp.IncrementEvent{Date: time.Now(), Kind: analyticsapp.IncrementAIUsage})
}

// IncrementWeKnoraUsage 增加 WeKnora 使用计数
func (s *StatisticsService) IncrementWeKnoraUsage(ctx context.Context) {
	_ = s.module.IncrementDailyStat(ctx, analyticsapp.IncrementEvent{Date: time.Now(), Kind: analyticsapp.IncrementWeKnora})
}

// StartDailyStatsWorker 启动每日统计后台任务
func (s *StatisticsService) StartDailyStatsWorker() {
	s.StartDailyStatsWorkerContext(context.Background(), 1*time.Hour)
}

// StartDailyStatsWorkerContext starts the stats worker with cancellation support.
func (s *StatisticsService) StartDailyStatsWorkerContext(ctx context.Context, interval time.Duration) {
	if interval <= 0 {
		interval = time.Hour
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	go func() {
		if err := s.UpdateDailyStats(context.Background(), time.Now()); err != nil {
			s.logger.Errorf("Failed to update daily stats: %v", err)
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := s.UpdateDailyStats(context.Background(), time.Now()); err != nil {
				s.logger.Errorf("Failed to update daily stats: %v", err)
			}
			yesterday := time.Now().AddDate(0, 0, -1)
			if err := s.UpdateDailyStats(context.Background(), yesterday); err != nil {
				s.logger.Errorf("Failed to update yesterday stats: %v", err)
			}
		}
	}
}

func dashboardStatsFromDTO(dto *analyticsapp.DashboardStats) *analyticscontract.DashboardStats {
	if dto == nil {
		return nil
	}
	return &analyticscontract.DashboardStats{
		TotalCustomers:       dto.TotalCustomers,
		TotalAgents:          dto.TotalAgents,
		TotalTickets:         dto.TotalTickets,
		TotalSessions:        dto.TotalSessions,
		TodayTickets:         dto.TodayTickets,
		TodaySessions:        dto.TodaySessions,
		TodayMessages:        dto.TodayMessages,
		OpenTickets:          dto.OpenTickets,
		AssignedTickets:      dto.AssignedTickets,
		ResolvedTickets:      dto.ResolvedTickets,
		ClosedTickets:        dto.ClosedTickets,
		OnlineAgents:         dto.OnlineAgents,
		BusyAgents:           dto.BusyAgents,
		ActiveSessions:       dto.ActiveSessions,
		AvgResponseTime:      dto.AvgResponseTime,
		AvgResolutionTime:    dto.AvgResolutionTime,
		CustomerSatisfaction: dto.CustomerSatisfaction,
		AIUsageToday:         dto.AIUsageToday,
		WeKnoraUsageToday:    dto.WeKnoraUsageToday,
	}
}

func timeRangeStatsFromDTO(items []analyticsapp.TimeRangeStats) []analyticscontract.TimeRangeStats {
	out := make([]analyticscontract.TimeRangeStats, 0, len(items))
	for _, item := range items {
		out = append(out, analyticscontract.TimeRangeStats{
			Date:                 item.Date,
			Tickets:              item.Tickets,
			Sessions:             item.Sessions,
			Messages:             item.Messages,
			ResolvedTickets:      item.ResolvedTickets,
			AvgResponseTime:      item.AvgResponseTime,
			CustomerSatisfaction: item.CustomerSatisfaction,
		})
	}
	return out
}

func agentPerformanceStatsFromDTO(items []analyticsapp.AgentPerformanceStats) []analyticscontract.AgentPerformanceStats {
	out := make([]analyticscontract.AgentPerformanceStats, 0, len(items))
	for _, item := range items {
		out = append(out, analyticscontract.AgentPerformanceStats{
			AgentID:           item.AgentID,
			AgentName:         item.AgentName,
			Department:        item.Department,
			TotalTickets:      item.TotalTickets,
			ResolvedTickets:   item.ResolvedTickets,
			AvgResponseTime:   item.AvgResponseTime,
			AvgResolutionTime: item.AvgResolutionTime,
			Rating:            item.Rating,
			OnlineTime:        item.OnlineTime,
		})
	}
	return out
}

func categoryStatsFromDTO(items []analyticsapp.CategoryStats) []analyticscontract.CategoryStats {
	out := make([]analyticscontract.CategoryStats, 0, len(items))
	for _, item := range items {
		out = append(out, analyticscontract.CategoryStats{
			Category: item.Category,
			Count:    item.Count,
		})
	}
	return out
}
