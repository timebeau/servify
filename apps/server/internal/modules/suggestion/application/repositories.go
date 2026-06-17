package application

import (
	"context"
	"time"
)

type TicketCandidate struct {
	ID          uint
	Title       string
	Description string
	Status      string
	Category    string
	Priority    string
	CreatedAt   time.Time
}

type KnowledgeDocCandidate struct {
	ID       uint
	Title    string
	Content  string
	Category string
	Tags     string
}

type Repository interface {
	FindTicketCandidates(ctx context.Context, tokens []string, candidateMax int) ([]TicketCandidate, error)
	FindKnowledgeDocCandidates(ctx context.Context, tokens []string) ([]KnowledgeDocCandidate, error)
}
