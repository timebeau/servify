package application

import (
	"context"

	"servify/apps/server/internal/modules/knowledge/domain"
)

type DocumentRepository interface {
	Create(ctx context.Context, doc *domain.Document) error
	Update(ctx context.Context, doc *domain.Document) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (*domain.Document, error)
	List(ctx context.Context, filter ListDocumentsFilter) ([]domain.Document, int64, error)
}

type IndexJobRepository interface {
	Create(ctx context.Context, job *domain.IndexJob) error
	Update(ctx context.Context, job *domain.IndexJob) error
	Get(ctx context.Context, id string) (*domain.IndexJob, error)
}
