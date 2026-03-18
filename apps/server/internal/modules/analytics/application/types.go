package application

import "time"

type DashboardStats struct {
	TotalCustomers       int64   `json:"total_customers"`
	TotalAgents          int64   `json:"total_agents"`
	TotalTickets         int64   `json:"total_tickets"`
	TotalSessions        int64   `json:"total_sessions"`
	TodayTickets         int64   `json:"today_tickets"`
	TodaySessions        int64   `json:"today_sessions"`
	TodayMessages        int64   `json:"today_messages"`
	OpenTickets          int64   `json:"open_tickets"`
	AssignedTickets      int64   `json:"assigned_tickets"`
	ResolvedTickets      int64   `json:"resolved_tickets"`
	ClosedTickets        int64   `json:"closed_tickets"`
	OnlineAgents         int64   `json:"online_agents"`
	BusyAgents           int64   `json:"busy_agents"`
	ActiveSessions       int64   `json:"active_sessions"`
	AvgResponseTime      float64 `json:"avg_response_time"`
	AvgResolutionTime    float64 `json:"avg_resolution_time"`
	CustomerSatisfaction float64 `json:"customer_satisfaction"`
	AIUsageToday         int64   `json:"ai_usage_today"`
	WeKnoraUsageToday    int64   `json:"weknora_usage_today"`
}

type TimeRangeStats struct {
	Date                 string  `json:"date"`
	Tickets              int64   `json:"tickets"`
	Sessions             int64   `json:"sessions"`
	Messages             int64   `json:"messages"`
	ResolvedTickets      int64   `json:"resolved_tickets"`
	AvgResponseTime      float64 `json:"avg_response_time"`
	CustomerSatisfaction float64 `json:"customer_satisfaction"`
}

type AgentPerformanceStats struct {
	AgentID           uint    `json:"agent_id"`
	AgentName         string  `json:"agent_name"`
	Department        string  `json:"department"`
	TotalTickets      int64   `json:"total_tickets"`
	ResolvedTickets   int64   `json:"resolved_tickets"`
	AvgResponseTime   float64 `json:"avg_response_time"`
	AvgResolutionTime float64 `json:"avg_resolution_time"`
	Rating            float64 `json:"rating"`
	OnlineTime        int64   `json:"online_time"`
}

type CategoryStats struct {
	Category string `json:"category"`
	Count    int64  `json:"count"`
}

type IncrementKind string

const (
	IncrementSessions IncrementKind = "sessions"
	IncrementMessages IncrementKind = "messages"
	IncrementTickets  IncrementKind = "tickets"
	IncrementResolved IncrementKind = "resolved"
	IncrementAIUsage  IncrementKind = "ai_usage"
	IncrementWeKnora  IncrementKind = "weknora_usage"
	IncrementSLA      IncrementKind = "sla_violations"
)

type IncrementEvent struct {
	Date time.Time
	Kind IncrementKind
}
