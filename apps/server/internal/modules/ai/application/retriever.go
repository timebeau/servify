package application

import (
	"context"
	"strings"

	"servify/apps/server/internal/platform/knowledgeprovider"
)

// Retriever isolates knowledge retrieval policy from the orchestrator.
type Retriever struct {
	provider knowledgeprovider.KnowledgeProvider
}

func NewRetriever(provider knowledgeprovider.KnowledgeProvider) *Retriever {
	return &Retriever{provider: provider}
}

func (r *Retriever) Retrieve(ctx context.Context, req AIRequest) ([]knowledgeprovider.KnowledgeHit, error) {
	if r == nil || r.provider == nil {
		return nil, nil
	}
	if !req.RetrievalPolicy.Enabled {
		return nil, nil
	}
	if strings.TrimSpace(req.Query) == "" {
		return nil, nil
	}

	searchReq := knowledgeprovider.SearchRequest{
		Query:           req.Query,
		TenantID:        req.TenantID,
		KnowledgeID:     req.ConversationID,
		TopK:            req.RetrievalPolicy.TopK,
		Threshold:       req.RetrievalPolicy.Threshold,
		Strategy:        req.RetrievalPolicy.Strategy,
		ConsistencyMode: knowledgeprovider.ConsistencyEventual,
	}
	return r.provider.Search(ctx, searchReq)
}
