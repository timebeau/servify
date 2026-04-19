package domain

import "time"

type Document struct {
	ID         string
	ProviderID string
	ExternalID string
	Title      string
	Content    string
	Category   string
	Tags       []string
	IsPublic   bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type IndexJobStatus string

const (
	IndexJobQueued  IndexJobStatus = "queued"
	IndexJobRunning IndexJobStatus = "running"
	IndexJobDone    IndexJobStatus = "done"
	IndexJobFailed  IndexJobStatus = "failed"
)

type IndexJob struct {
	ID          string
	DocumentID  string
	Status      IndexJobStatus
	Error       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CompletedAt *time.Time
}
