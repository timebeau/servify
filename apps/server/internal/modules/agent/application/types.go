package application

import "time"

type CreateAgentCommand struct {
	UserID             uint
	Department         string
	Skills             []string
	MaxChatConcurrency int
}

type AgentRuntimeDTO struct {
	UserID              uint      `json:"user_id"`
	Username            string    `json:"username"`
	Name                string    `json:"name"`
	Department          string    `json:"department"`
	Skills              []string  `json:"skills"`
	Status              string    `json:"status"`
	MaxChatConcurrency  int       `json:"max_chat_concurrency"`
	MaxVoiceConcurrency int       `json:"max_voice_concurrency"`
	CurrentChatLoad     int       `json:"current_chat_load"`
	CurrentVoiceLoad    int       `json:"current_voice_load"`
	Rating              float64   `json:"rating"`
	AvgResponseTime     int       `json:"avg_response_time"`
	LastActivity        time.Time `json:"last_activity"`
	ConnectedAt         time.Time `json:"connected_at"`
}

type AgentStatsDTO struct {
	Total           int64   `json:"total"`
	Online          int64   `json:"online"`
	Busy            int64   `json:"busy"`
	AvgResponseTime int64   `json:"avg_response_time"`
	AvgRating       float64 `json:"avg_rating"`
}
