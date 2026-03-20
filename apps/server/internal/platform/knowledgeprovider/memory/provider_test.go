package memory

import (
	"context"
	"testing"

	"servify/apps/server/internal/platform/knowledgeprovider"
)

func TestProviderSearchAndNamespaceMapping(t *testing.T) {
	provider := NewProvider("tenant-default", "kb-default")
	if err := provider.UpsertDocument(context.Background(), knowledgeprovider.KnowledgeDocument{
		ID:       "doc-1",
		TenantID: "tenant-a",
		Title:    "Billing",
		Content:  "Billing details",
	}); err != nil {
		t.Fatalf("upsert: %v", err)
	}

	hits, err := provider.Search(context.Background(), knowledgeprovider.SearchRequest{
		Query:    "billing",
		TenantID: "tenant-a",
	})
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(hits) != 1 {
		t.Fatalf("expected 1 hit, got %d", len(hits))
	}
	if hits[0].Source != "memory" {
		t.Fatalf("unexpected source: %+v", hits[0])
	}
}
