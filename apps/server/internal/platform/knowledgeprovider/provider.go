package knowledgeprovider

import (
	"context"
	"errors"
)

var ErrOperationNotSupported = errors.New("knowledge provider operation is not supported")

// KnowledgeProvider defines the contract for retrieval and indexing providers.
type KnowledgeProvider interface {
	Search(ctx context.Context, req SearchRequest) ([]KnowledgeHit, error)
	UpsertDocument(ctx context.Context, doc KnowledgeDocument) (string, error)
	DeleteDocument(ctx context.Context, id string) error
	HealthCheck(ctx context.Context) error
}

// RebuildableProvider is an optional extension for providers that can rebuild an index.
type RebuildableProvider interface {
	RebuildIndex(ctx context.Context, req RebuildRequest) error
}
