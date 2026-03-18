package domain

type DashboardReadModel struct {
	TotalCustomers       int64
	TotalAgents          int64
	TotalTickets         int64
	TotalSessions        int64
	TodayTickets         int64
	TodaySessions        int64
	TodayMessages        int64
	OpenTickets          int64
	AssignedTickets      int64
	ResolvedTickets      int64
	ClosedTickets        int64
	OnlineAgents         int64
	BusyAgents           int64
	ActiveSessions       int64
	AvgResponseTime      float64
	AvgResolutionTime    float64
	CustomerSatisfaction float64
	AIUsageToday         int64
	WeKnoraUsageToday    int64
}

type TicketTrendReadModel struct {
	Date                 string
	Tickets              int64
	Sessions             int64
	Messages             int64
	ResolvedTickets      int64
	AvgResponseTime      float64
	CustomerSatisfaction float64
}

type AgentPerformanceReadModel struct {
	AgentID           uint
	AgentName         string
	Department        string
	TotalTickets      int64
	ResolvedTickets   int64
	AvgResponseTime   float64
	AvgResolutionTime float64
	Rating            float64
	OnlineTime        int64
}

type SatisfactionTrendReadModel struct {
	Date  string
	Score float64
}

type SLATrendReadModel struct {
	Date       string
	Violations int64
}
