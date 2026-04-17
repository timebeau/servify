package application

import "time"

type CreateDocumentRequest struct {
	ID       string
	Title    string
	Content  string
	Category string
	Tags     []string
	IsPublic bool
}

type UpdateDocumentRequest struct {
	Title    *string
	Content  *string
	Category *string
	Tags     *[]string
	IsPublic *bool
}

type ListDocumentsFilter struct {
	Page       int
	PageSize   int
	Category   string
	Search     string
	PublicOnly bool
}

type QueueIndexJobRequest struct {
	JobID      string
	DocumentID string
}

type RunIndexJobRequest struct {
	JobID string
}

type IndexJobResult struct {
	JobID       string
	DocumentID  string
	Status      string
	Error       string
	CompletedAt *time.Time
}
