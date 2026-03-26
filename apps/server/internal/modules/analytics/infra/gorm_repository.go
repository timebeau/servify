package infra

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"servify/apps/server/internal/models"
	analyticsapp "servify/apps/server/internal/modules/analytics/application"
	platformauth "servify/apps/server/internal/platform/auth"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) GetDashboardStats(ctx context.Context) (*analyticsapp.DashboardStats, error) {
	stats := &analyticsapp.DashboardStats{}
	today := time.Now().Truncate(24 * time.Hour)
	customerScope(r.db.WithContext(ctx).Model(&models.User{}), ctx).Where("role = ?", "customer").Count(&stats.TotalCustomers)
	applyEntityScope(r.db.WithContext(ctx).Model(&models.Agent{}), ctx).Count(&stats.TotalAgents)
	applyEntityScope(r.db.WithContext(ctx).Model(&models.Ticket{}), ctx).Count(&stats.TotalTickets)
	applyEntityScope(r.db.WithContext(ctx).Model(&models.Session{}), ctx).Count(&stats.TotalSessions)
	applyEntityScope(r.db.WithContext(ctx).Model(&models.Ticket{}), ctx).Where("created_at >= ?", today).Count(&stats.TodayTickets)
	applyEntityScope(r.db.WithContext(ctx).Model(&models.Session{}), ctx).Where("created_at >= ?", today).Count(&stats.TodaySessions)
	applyEntityScope(r.db.WithContext(ctx).Model(&models.Message{}), ctx).Where("created_at >= ?", today).Count(&stats.TodayMessages)
	applyEntityScope(r.db.WithContext(ctx).Model(&models.Ticket{}), ctx).Where("status = ?", "open").Count(&stats.OpenTickets)
	applyEntityScope(r.db.WithContext(ctx).Model(&models.Ticket{}), ctx).Where("status = ?", "assigned").Count(&stats.AssignedTickets)
	applyEntityScope(r.db.WithContext(ctx).Model(&models.Ticket{}), ctx).Where("status = ?", "resolved").Count(&stats.ResolvedTickets)
	applyEntityScope(r.db.WithContext(ctx).Model(&models.Ticket{}), ctx).Where("status = ?", "closed").Count(&stats.ClosedTickets)
	applyEntityScope(r.db.WithContext(ctx).Model(&models.Agent{}), ctx).Where("status = ?", "online").Count(&stats.OnlineAgents)
	applyEntityScope(r.db.WithContext(ctx).Model(&models.Agent{}), ctx).Where("status = ?", "busy").Count(&stats.BusyAgents)
	applyEntityScope(r.db.WithContext(ctx).Model(&models.Session{}), ctx).Where("status = ?", "active").Count(&stats.ActiveSessions)
	applyEntityScope(r.db.WithContext(ctx).Model(&models.Agent{}), ctx).Select("AVG(avg_response_time)").Row().Scan(&stats.AvgResponseTime)
	var avgResolution float64
	applyEntityScope(r.db.WithContext(ctx).Model(&models.Ticket{}), ctx).Where("resolved_at IS NOT NULL").Select("AVG(EXTRACT(epoch FROM (resolved_at - created_at)))").Row().Scan(&avgResolution)
	stats.AvgResolutionTime = avgResolution
	stats.CustomerSatisfaction = 4.2
	var dailyStat models.DailyStats
	if err := r.db.WithContext(ctx).Where("date = ?", today).First(&dailyStat).Error; err == nil {
		stats.AIUsageToday = int64(dailyStat.AIUsageCount)
		stats.WeKnoraUsageToday = int64(dailyStat.WeKnoraUsageCount)
	}
	return stats, nil
}

func (r *GormRepository) GetTimeRangeStats(ctx context.Context, startDate, endDate time.Time) ([]analyticsapp.TimeRangeStats, error) {
	var stats []analyticsapp.TimeRangeStats
	current := startDate.Truncate(24 * time.Hour)
	end := endDate.Truncate(24 * time.Hour)
	for current.Before(end) || current.Equal(end) {
		nextDay := current.Add(24 * time.Hour)
		stat := analyticsapp.TimeRangeStats{Date: current.Format("2006-01-02")}
		applyEntityScope(r.db.WithContext(ctx).Model(&models.Ticket{}), ctx).Where("created_at >= ? AND created_at < ?", current, nextDay).Count(&stat.Tickets)
		applyEntityScope(r.db.WithContext(ctx).Model(&models.Session{}), ctx).Where("created_at >= ? AND created_at < ?", current, nextDay).Count(&stat.Sessions)
		applyEntityScope(r.db.WithContext(ctx).Model(&models.Message{}), ctx).Where("created_at >= ? AND created_at < ?", current, nextDay).Count(&stat.Messages)
		applyEntityScope(r.db.WithContext(ctx).Model(&models.Ticket{}), ctx).Where("resolved_at >= ? AND resolved_at < ?", current, nextDay).Count(&stat.ResolvedTickets)
		var daily models.DailyStats
		if err := r.db.WithContext(ctx).Where("date = ?", current).First(&daily).Error; err == nil {
			stat.AvgResponseTime = float64(daily.AvgResponseTime)
			stat.CustomerSatisfaction = daily.CustomerSatisfaction
		}
		stats = append(stats, stat)
		current = nextDay
	}
	return stats, nil
}

func (r *GormRepository) GetAgentPerformanceStats(ctx context.Context, startDate, endDate time.Time, limit int) ([]analyticsapp.AgentPerformanceStats, error) {
	var stats []analyticsapp.AgentPerformanceStats
	query := `
		SELECT
			a.user_id as agent_id,
			u.name as agent_name,
			a.department,
			COUNT(t.id) as total_tickets,
			COUNT(CASE WHEN t.status = 'resolved' OR t.status = 'closed' THEN 1 END) as resolved_tickets,
			a.avg_response_time,
			AVG(CASE WHEN t.resolved_at IS NOT NULL
				THEN EXTRACT(epoch FROM (t.resolved_at - t.created_at))
				END) as avg_resolution_time,
			a.rating
		FROM agents a
		LEFT JOIN users u ON a.user_id = u.id
		LEFT JOIN tickets t ON a.user_id = t.agent_id
			AND t.created_at >= ? AND t.created_at <= ?
			AND (? = '' OR t.tenant_id = ?)
			AND (? = '' OR t.workspace_id = ?)
		WHERE (? = '' OR a.tenant_id = ?)
			AND (? = '' OR a.workspace_id = ?)
		GROUP BY a.user_id, u.name, a.department, a.avg_response_time, a.rating
		ORDER BY total_tickets DESC
	`
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}
	tenantID, workspaceID := scopeValues(ctx)
	if err := r.db.WithContext(ctx).Raw(query, startDate, endDate, tenantID, tenantID, workspaceID, workspaceID, tenantID, tenantID, workspaceID, workspaceID).Scan(&stats).Error; err != nil {
		return nil, fmt.Errorf("failed to get agent performance stats: %w", err)
	}
	return stats, nil
}

func (r *GormRepository) GetTicketCategoryStats(ctx context.Context, startDate, endDate time.Time) ([]analyticsapp.CategoryStats, error) {
	var stats []analyticsapp.CategoryStats
	err := applyEntityScope(r.db.WithContext(ctx).Model(&models.Ticket{}), ctx).Select("category, COUNT(*) as count").Where("created_at >= ? AND created_at <= ?", startDate, endDate).Group("category").Order("count DESC").Scan(&stats).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get category stats: %w", err)
	}
	return stats, nil
}

func (r *GormRepository) GetTicketPriorityStats(ctx context.Context, startDate, endDate time.Time) ([]analyticsapp.CategoryStats, error) {
	var stats []analyticsapp.CategoryStats
	err := applyEntityScope(r.db.WithContext(ctx).Model(&models.Ticket{}), ctx).Select("priority as category, COUNT(*) as count").Where("created_at >= ? AND created_at <= ?", startDate, endDate).Group("priority").Order("count DESC").Scan(&stats).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get priority stats: %w", err)
	}
	return stats, nil
}

func (r *GormRepository) GetCustomerSourceStats(ctx context.Context) ([]analyticsapp.CategoryStats, error) {
	var stats []analyticsapp.CategoryStats
	err := applyEntityScope(r.db.WithContext(ctx).Model(&models.Customer{}), ctx).Select("source as category, COUNT(*) as count").Group("source").Order("count DESC").Scan(&stats).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get customer source stats: %w", err)
	}
	return stats, nil
}

func (r *GormRepository) UpdateDailyStats(ctx context.Context, date time.Time) error {
	date = date.Truncate(24 * time.Hour)
	nextDay := date.Add(24 * time.Hour)
	var dailyStats models.DailyStats
	err := r.db.WithContext(ctx).Where("date = ?", date).First(&dailyStats).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			dailyStats = models.DailyStats{Date: date}
		} else {
			return fmt.Errorf("failed to query daily stats: %w", err)
		}
	}
	var totalSessions, totalMessages, totalTickets, resolvedTickets int64
	applyEntityScope(r.db.WithContext(ctx).Model(&models.Session{}), ctx).Where("created_at >= ? AND created_at < ?", date, nextDay).Count(&totalSessions)
	applyEntityScope(r.db.WithContext(ctx).Model(&models.Message{}), ctx).Where("created_at >= ? AND created_at < ?", date, nextDay).Count(&totalMessages)
	applyEntityScope(r.db.WithContext(ctx).Model(&models.Ticket{}), ctx).Where("created_at >= ? AND created_at < ?", date, nextDay).Count(&totalTickets)
	applyEntityScope(r.db.WithContext(ctx).Model(&models.Ticket{}), ctx).Where("resolved_at >= ? AND resolved_at < ?", date, nextDay).Count(&resolvedTickets)
	dailyStats.TotalSessions = int(totalSessions)
	dailyStats.TotalMessages = int(totalMessages)
	dailyStats.TotalTickets = int(totalTickets)
	dailyStats.ResolvedTickets = int(resolvedTickets)
	var avgResponseTime, avgResolutionTime float64
	applyEntityScope(r.db.WithContext(ctx).Model(&models.Agent{}), ctx).Select("AVG(avg_response_time)").Row().Scan(&avgResponseTime)
	applyEntityScope(r.db.WithContext(ctx).Model(&models.Ticket{}), ctx).Where("resolved_at >= ? AND resolved_at < ? AND resolved_at IS NOT NULL", date, nextDay).Select("AVG(EXTRACT(epoch FROM (resolved_at - created_at)))").Row().Scan(&avgResolutionTime)
	dailyStats.AvgResponseTime = int(avgResponseTime)
	dailyStats.AvgResolutionTime = int(avgResolutionTime)
	dailyStats.CustomerSatisfaction = 4.2
	if dailyStats.ID == 0 {
		err = r.db.WithContext(ctx).Create(&dailyStats).Error
	} else {
		err = r.db.WithContext(ctx).Save(&dailyStats).Error
	}
	if err != nil {
		return fmt.Errorf("failed to save daily stats: %w", err)
	}
	return nil
}

func scopeValues(ctx context.Context) (string, string) {
	return platformauth.TenantIDFromContext(ctx), platformauth.WorkspaceIDFromContext(ctx)
}

func applyEntityScope(db *gorm.DB, ctx context.Context) *gorm.DB {
	tenantID, workspaceID := scopeValues(ctx)
	if tenantID != "" {
		db = db.Where("tenant_id = ?", tenantID)
	}
	if workspaceID != "" {
		db = db.Where("workspace_id = ?", workspaceID)
	}
	return db
}

func customerScope(db *gorm.DB, ctx context.Context) *gorm.DB {
	tenantID, workspaceID := scopeValues(ctx)
	if tenantID == "" && workspaceID == "" {
		return db
	}
	db = db.Joins("JOIN customers ON customers.user_id = users.id")
	if tenantID != "" {
		db = db.Where("customers.tenant_id = ?", tenantID)
	}
	if workspaceID != "" {
		db = db.Where("customers.workspace_id = ?", workspaceID)
	}
	return db
}

func (r *GormRepository) IncrementDailyStat(ctx context.Context, event analyticsapp.IncrementEvent) error {
	date := event.Date.Truncate(24 * time.Hour)
	column := ""
	switch event.Kind {
	case analyticsapp.IncrementSessions:
		column = "total_sessions"
	case analyticsapp.IncrementMessages:
		column = "total_messages"
	case analyticsapp.IncrementTickets:
		column = "total_tickets"
	case analyticsapp.IncrementResolved:
		column = "resolved_tickets"
	case analyticsapp.IncrementAIUsage:
		column = "ai_usage_count"
	case analyticsapp.IncrementWeKnora:
		column = "we_knora_usage_count"
	case analyticsapp.IncrementSLA:
		column = "sla_violations"
	default:
		return nil
	}
	var daily models.DailyStats
	if err := r.db.WithContext(ctx).Where("date = ?", date).First(&daily).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			daily = models.DailyStats{Date: date}
			if err := r.db.WithContext(ctx).Create(&daily).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return r.db.WithContext(ctx).Model(&daily).UpdateColumn(column, gorm.Expr(column+" + 1")).Error
}
