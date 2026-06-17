package contract

type Badge struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type LeaderboardEntry struct {
	Rank            int     `json:"rank"`
	AgentID         uint    `json:"agent_id"`
	AgentName       string  `json:"agent_name"`
	Username        string  `json:"username"`
	Department      string  `json:"department"`
	ResolvedTickets int64   `json:"resolved_tickets"`
	CSATAvg         float64 `json:"csat_avg"`
	CSATCount       int64   `json:"csat_count"`
	AvgResponseTime int64   `json:"avg_response_time"`
	Score           float64 `json:"score"`
	Badges          []Badge `json:"badges,omitempty"`
}

type LeaderboardResponse struct {
	StartDate string             `json:"start_date"`
	EndDate   string             `json:"end_date"`
	Limit     int                `json:"limit"`
	Entries   []LeaderboardEntry `json:"entries"`
}
