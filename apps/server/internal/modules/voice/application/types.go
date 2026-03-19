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

type HoldCallCommand struct {
	CallID string
}

type ResumeCallCommand struct {
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
	ID              string     `json:"id"`
	SessionID       string     `json:"session_id"`
	Status          string     `json:"status"`
	StartedAt       time.Time  `json:"started_at"`
	AnsweredAt      *time.Time `json:"answered_at,omitempty"`
	HeldAt          *time.Time `json:"held_at,omitempty"`
	ResumedAt       *time.Time `json:"resumed_at,omitempty"`
	EndedAt         *time.Time `json:"ended_at,omitempty"`
	TransferToAgent *uint      `json:"transfer_to_agent,omitempty"`
}

type RecordingDTO struct {
	ID        string    `json:"id"`
	CallID    string    `json:"call_id"`
	Provider  string    `json:"provider,omitempty"`
	Status    string    `json:"status"`
	StartedAt time.Time `json:"started_at"`
}

type TranscriptDTO struct {
	CallID     string    `json:"call_id"`
	Content    string    `json:"content"`
	Language   string    `json:"language"`
	Finalized  bool      `json:"finalized"`
	AppendedAt time.Time `json:"appended_at"`
}
