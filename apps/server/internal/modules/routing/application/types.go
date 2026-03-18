package application

import (
	"time"

	"servify/apps/server/internal/modules/routing/domain"
)

type RequestHumanHandoffCommand struct {
	SessionID    string
	Reason       string
	TargetSkills []string
	Priority     string
	Notes        string
}

type AssignAgentCommand struct {
	SessionID   string
	AgentID     uint
	FromAgentID *uint
	Reason      string
	Notes       string
}

type AddToWaitingQueueCommand struct {
	SessionID    string
	Reason       string
	TargetSkills []string
	Priority     string
	Notes        string
}

type CancelWaitingCommand struct {
	SessionID string
	Reason    string
}

type AssignmentDTO struct {
	SessionID   string    `json:"session_id"`
	FromAgentID *uint     `json:"from_agent_id,omitempty"`
	ToAgentID   uint      `json:"to_agent_id"`
	Reason      string    `json:"reason,omitempty"`
	Notes       string    `json:"notes,omitempty"`
	AssignedAt  time.Time `json:"assigned_at"`
}

type QueueEntryDTO struct {
	SessionID    string     `json:"session_id"`
	Reason       string     `json:"reason,omitempty"`
	TargetSkills []string   `json:"target_skills,omitempty"`
	Priority     string     `json:"priority,omitempty"`
	Notes        string     `json:"notes,omitempty"`
	Status       string     `json:"status"`
	QueuedAt     time.Time  `json:"queued_at"`
	AssignedAt   *time.Time `json:"assigned_at,omitempty"`
	AssignedTo   *uint      `json:"assigned_to,omitempty"`
}

func MapAssignment(item domain.Assignment) AssignmentDTO {
	return AssignmentDTO{
		SessionID:   item.SessionID,
		FromAgentID: item.FromAgentID,
		ToAgentID:   item.ToAgentID,
		Reason:      item.Reason,
		Notes:       item.Notes,
		AssignedAt:  item.AssignedAt,
	}
}

func MapQueueEntry(item domain.QueueEntry) QueueEntryDTO {
	return QueueEntryDTO{
		SessionID:    item.SessionID,
		Reason:       item.Reason,
		TargetSkills: append([]string(nil), item.TargetSkills...),
		Priority:     item.Priority,
		Notes:        item.Notes,
		Status:       string(item.Status),
		QueuedAt:     item.QueuedAt,
		AssignedAt:   item.AssignedAt,
		AssignedTo:   item.AssignedTo,
	}
}
