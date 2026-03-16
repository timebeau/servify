package weknora

import (
	"context"
	"testing"
	"time"

	"servify/apps/server/internal/platform/knowledgeprovider"
	base "servify/apps/server/pkg/weknora"
)

type mockClient struct {
	searchResp *base.SearchResponse
	searchErr  error
	uploadErr  error
	healthyErr error
}

func (m *mockClient) CreateKnowledgeBase(ctx context.Context, req *base.CreateKBRequest) (*base.KnowledgeBase, error) {
	return nil, nil
}
func (m *mockClient) GetKnowledgeBase(ctx context.Context, kbID string) (*base.KnowledgeBase, error) {
	return nil, nil
}
func (m *mockClient) UploadDocument(ctx context.Context, kbID string, doc *base.Document) (*base.DocumentInfo, error) {
	return &base.DocumentInfo{ID: "doc-1", Title: doc.Title, ProcessedAt: time.Now()}, m.uploadErr
}
func (m *mockClient) SearchKnowledge(ctx context.Context, req *base.SearchRequest) (*base.SearchResponse, error) {
	return m.searchResp, m.searchErr
}
func (m *mockClient) CreateSession(ctx context.Context, req *base.SessionRequest) (*base.Session, error) {
	return nil, nil
}
func (m *mockClient) Chat(ctx context.Context, sessionID string, req *base.ChatRequest) (*base.ChatResponse, error) {
	return nil, nil
}
func (m *mockClient) HealthCheck(ctx context.Context) error { return m.healthyErr }

func TestProviderSearch(t *testing.T) {
	provider := NewProvider(&mockClient{
		searchResp: &base.SearchResponse{
			Success: true,
			Data: base.SearchData{
				Results: []base.SearchResult{
					{DocumentID: "doc-1", Title: "Billing", Content: "Billing content", Score: 0.95, Source: "weknora"},
				},
			},
		},
	}, "kb-1")

	hits, err := provider.Search(context.Background(), knowledgeprovider.SearchRequest{
		Query: "billing",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(hits) != 1 {
		t.Fatalf("expected 1 hit, got %d", len(hits))
	}
	if hits[0].DocumentID != "doc-1" {
		t.Fatalf("unexpected document id: %s", hits[0].DocumentID)
	}
}

func TestProviderHealthCheck(t *testing.T) {
	provider := NewProvider(&mockClient{}, "kb-1")
	if err := provider.HealthCheck(context.Background()); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
