package application

import "time"

type LeaderboardRequest struct {
	StartDate  time.Time
	EndDate    time.Time
	Limit      int
	Department string
}

type Badge struct {
	ID          string
	Name        string
	Description string
}

type LeaderboardEntry struct {
	Rank            int
	AgentID         uint
	AgentName       string
	Username        string
	Department      string
	ResolvedTickets int64
	CSATAvg         float64
	CSATCount       int64
	AvgResponseTime int64
	Score           float64
	Badges          []Badge
}

type LeaderboardResponse struct {
	StartDate string
	EndDate   string
	Limit     int
	Entries   []LeaderboardEntry
}

type AgentProfile struct {
	UserID          uint
	Username        string
	Name            string
	Department      string
	AvgResponseTime int64
}

type AgentResolvedCount struct {
	AgentID uint
	Count   int64
}

type AgentCSAT struct {
	AgentID uint
	Avg     float64
	Count   int64
}
