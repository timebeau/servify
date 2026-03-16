package knowledgeprovider

import "context"

// KnowledgeProvider defines the contract for retrieval and indexing providers.
type KnowledgeProvider interface {
	Search(ctx context.Context, req SearchRequest) ([]KnowledgeHit, error)
	UpsertDocument(ctx context.Context, doc KnowledgeDocument) error
	DeleteDocument(ctx context.Context, id string) error
	HealthCheck(ctx context.Context) error
}
