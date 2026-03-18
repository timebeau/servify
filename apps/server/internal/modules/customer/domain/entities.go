package domain

import "time"

type CustomerProfile struct {
	UserID    uint
	Username  string
	Email     string
	Name      string
	Phone     string
	Status    string
	Company   string
	Industry  string
	Source    string
	Priority  string
	Tags      []CustomerTag
	Notes     []CustomerNote
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CustomerTag struct {
	Value string
}

type CustomerNote struct {
	AuthorID  uint
	Content   string
	CreatedAt time.Time
}

type CustomerActivity struct {
	CustomerID uint
}
