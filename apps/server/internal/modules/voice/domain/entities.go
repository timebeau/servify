package domain

import "time"

type CallStatus string

const (
	CallStatusStarted  CallStatus = "started"
	CallStatusAnswered CallStatus = "answered"
	CallStatusEnded    CallStatus = "ended"
)

type CallSession struct {
	ID           string
	SessionID    string
	Status       CallStatus
	StartedAt    time.Time
	AnsweredAt   *time.Time
	EndedAt      *time.Time
	Participants []VoiceParticipant
}

type MediaSession struct {
	ID         string
	CallID     string
	Kind       string
	Status     string
	ExternalID string
	CreatedAt  time.Time
}

type VoiceParticipant struct {
	ID       string
	UserID   *uint
	Role     string
	Muted    bool
	JoinedAt time.Time
}

type Recording struct {
	ID         string
	CallID     string
	Status     string
	StorageURI string
	CreatedAt  time.Time
}

type Transcript struct {
	ID        string
	CallID    string
	Content   string
	Language  string
	CreatedAt time.Time
}
