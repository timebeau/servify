package domain

import "time"

type Ticket struct {
	ID          uint
	Title       string
	Description string
	CustomerID  uint
	AgentID     *uint
	SessionID   *string
	Category    string
	Priority    string
	Status      string
	Source      string
	Tags        string
	DueDate     *time.Time
	ResolvedAt  *time.Time
	ClosedAt    *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CustomFieldValue struct {
	CustomFieldID uint
	Key           string
	Value         string
}

type Comment struct {
	ID        uint
	UserID    uint
	Content   string
	Type      string
	CreatedAt time.Time
}

type StatusChange struct {
	ID         uint
	UserID     uint
	FromStatus string
	ToStatus   string
	Reason     string
	CreatedAt  time.Time
}

type TicketDetails struct {
	Ticket
	CustomFieldValues []CustomFieldValue
	Comments          []Comment
	StatusHistory     []StatusChange
}
