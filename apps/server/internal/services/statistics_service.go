package services

import (
	"context"
	"time"

	analyticsapp "servify/apps/server/internal/modules/analytics/application"
	analyticsdelivery "servify/apps/server/internal/modules/analytics/delivery"
	analyticsinfra "servify/apps/server/internal/modules/analytics/infra"
	"servify/apps/server/internal/platform/eventbus"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// StatisticsService 数据统计服务兼容层。
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

// DashboardStats 仪表板统计数据
type DashboardStats = analyticsapp.DashboardStats

// TimeRangeStats 时间范围统计
type TimeRangeStats = analyticsapp.TimeRangeStats

// AgentPerformanceStats 客服绩效统计
type AgentPerformanceStats = analyticsapp.AgentPerformanceStats

// CategoryStats 分类统计
type CategoryStats = analyticsapp.CategoryStats

// GetDashboardStats 获取仪表板统计数据
func (s *StatisticsService) GetDashboardStats(ctx context.Context) (*DashboardStats, error) {
	return s.module.GetDashboardStats(ctx)
}

// GetTimeRangeStats 获取时间范围统计
func (s *StatisticsService) GetTimeRangeStats(ctx context.Context, startDate, endDate time.Time) ([]TimeRangeStats, error) {
	return s.module.GetTimeRangeStats(ctx, startDate, endDate)
}

// GetAgentPerformanceStats 获取客服绩效统计
func (s *StatisticsService) GetAgentPerformanceStats(ctx context.Context, startDate, endDate time.Time, limit int) ([]AgentPerformanceStats, error) {
	return s.module.GetAgentPerformanceStats(ctx, startDate, endDate, limit)
}

// GetTicketCategoryStats 获取工单分类统计
func (s *StatisticsService) GetTicketCategoryStats(ctx context.Context, startDate, endDate time.Time) ([]CategoryStats, error) {
	return s.module.GetTicketCategoryStats(ctx, startDate, endDate)
}

// GetTicketPriorityStats 获取工单优先级统计
func (s *StatisticsService) GetTicketPriorityStats(ctx context.Context, startDate, endDate time.Time) ([]CategoryStats, error) {
	return s.module.GetTicketPriorityStats(ctx, startDate, endDate)
}

// GetCustomerSourceStats 获取客户来源统计
func (s *StatisticsService) GetCustomerSourceStats(ctx context.Context) ([]CategoryStats, error) {
	return s.module.GetCustomerSourceStats(ctx)
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
