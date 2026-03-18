package application

import (
	"time"

	"servify/apps/server/internal/models"
)

type CreateCustomerCommand struct {
	Username string
	Email    string
	Name     string
	Phone    string
	Company  string
	Industry string
	Source   string
	Tags     []string
	Notes    string
	Priority string
}

type UpdateCustomerCommand struct {
	Name     *string
	Phone    *string
	Company  *string
	Industry *string
	Source   *string
	Tags     *[]string
	Notes    *string
	Priority *string
	Status   *string
}

type ListCustomersQuery struct {
	Page      int
	PageSize  int
	Search    string
	Industry  []string
	Source    []string
	Priority  []string
	Status    []string
	Tags      []string
	SortBy    string
	SortOrder string
}

type CustomerInfoDTO struct {
	models.User
	Company  string `json:"company"`
	Industry string `json:"industry"`
	Source   string `json:"source"`
	Tags     string `json:"tags"`
	Notes    string `json:"notes"`
	Priority string `json:"priority"`
}

type CustomerActivityDTO struct {
	CustomerID     uint             `json:"customer_id"`
	RecentSessions []models.Session `json:"recent_sessions"`
	RecentTickets  []models.Ticket  `json:"recent_tickets"`
	RecentMessages []models.Message `json:"recent_messages"`
}

type CustomerStatsDTO struct {
	Total       int64           `json:"total"`
	Active      int64           `json:"active"`
	NewThisWeek int64           `json:"new_this_week"`
	BySource    []SourceCount   `json:"by_source"`
	ByIndustry  []IndustryCount `json:"by_industry"`
	ByPriority  []PriorityCount `json:"by_priority"`
}

type SourceCount struct {
	Source string `json:"source"`
	Count  int64  `json:"count"`
}

type IndustryCount struct {
	Industry string `json:"industry"`
	Count    int64  `json:"count"`
}

type PriorityCount struct {
	Priority string `json:"priority"`
	Count    int64  `json:"count"`
}

type CustomerNoteDTO struct {
	AuthorID  uint      `json:"author_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
