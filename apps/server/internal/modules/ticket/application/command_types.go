package application

import "time"

type CreateTicketCommand struct {
	Title       string
	Description string
	CustomerID  uint
	Category    string
	Priority    string
	Source      string
	Tags        string
	SessionID   *string
}

type UpdateTicketCommand struct {
	Title       *string
	Description *string
	AgentID     *uint
	Category    *string
	Priority    *string
	Status      *string
	Tags        *string
	DueDate     *time.Time
	ActorID     uint
}

type AssignTicketCommand struct {
	AgentID uint
	UserID  uint
}

type AddCommentCommand struct {
	UserID      uint
	Content     string
	CommentType string
}

type CloseTicketCommand struct {
	UserID uint
	Reason string
}
