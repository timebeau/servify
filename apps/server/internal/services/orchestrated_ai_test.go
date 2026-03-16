package services

import (
	"context"
	"testing"

	mockkp "servify/apps/server/internal/platform/knowledgeprovider/mock"
	"servify/apps/server/internal/platform/knowledgeprovider"
	mockllm "servify/apps/server/internal/platform/llm/mock"
	"servify/apps/server/internal/platform/llm"
)

func TestOrchestratedAIServiceProcessQuery(t *testing.T) {
	svc := NewOrchestratedAIService(
		&mockllm.Provider{
			ChatResponse: llm.ChatResponse{Content: "orchestrated answer"},
		},
		&mockkp.Provider{
			Hits: []knowledgeprovider.KnowledgeHit{
				{DocumentID: "doc-1", Title: "Billing", Content: "Billing details", Score: 0.9},
			},
		},
	)

	resp, err := svc.ProcessQuery(context.Background(), "billing", "session-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Content != "orchestrated answer" {
		t.Fatalf("unexpected content: %s", resp.Content)
	}
	if resp.Confidence <= 0.6 {
		t.Fatalf("expected higher confidence with sources, got %v", resp.Confidence)
	}
}

func TestOrchestratedAIServiceStatus(t *testing.T) {
	svc := NewOrchestratedAIService(&mockllm.Provider{}, &mockkp.Provider{})
	st := svc.GetStatus(context.Background())
	if st["type"] != "orchestrated" {
		t.Fatalf("unexpected type: %v", st["type"])
	}
}
