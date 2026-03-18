package application

import "time"

type StartCallCommand struct {
	CallID       string
	SessionID    string
	ConnectionID string
}

type AnswerCallCommand struct {
	CallID string
}

type EndCallCommand struct {
	CallID string
}

type TransferCallCommand struct {
	CallID    string
	ToAgentID uint
}

type CallDTO struct {
	ID         string     `json:"id"`
	SessionID  string     `json:"session_id"`
	Status     string     `json:"status"`
	StartedAt  time.Time  `json:"started_at"`
	AnsweredAt *time.Time `json:"answered_at,omitempty"`
	EndedAt    *time.Time `json:"ended_at,omitempty"`
}
