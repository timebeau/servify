package application

import (
	"context"
	"testing"

	mockkp "servify/apps/server/internal/platform/knowledgeprovider/mock"
	"servify/apps/server/internal/platform/knowledgeprovider"
	mockllm "servify/apps/server/internal/platform/llm/mock"
	"servify/apps/server/internal/platform/llm"
)

func TestQueryOrchestratorHandleWithRetrieval(t *testing.T) {
	llmProvider := &mockllm.Provider{
		ChatResponse: llm.ChatResponse{
			Content:      "answer",
			Model:        "mock-model",
			FinishReason: "stop",
			TokenUsage: &llm.TokenUsage{
				InputTokens:  10,
				OutputTokens: 20,
				TotalTokens:  30,
			},
		},
	}
	knowledgeProvider := &mockkp.Provider{
		Hits: []knowledgeprovider.KnowledgeHit{
			{DocumentID: "doc-1", Title: "Billing", Content: "Billing help content", Score: 0.9},
		},
	}

	orchestrator := NewQueryOrchestrator(llmProvider, knowledgeProvider)
	resp, err := orchestrator.Handle(context.Background(), AIRequest{
		TenantID:     "tenant-1",
		TaskType:     TaskTypeQA,
		Query:        "billing",
		SystemPrompt: "You are helpful.",
		RetrievalPolicy: RetrievalPolicy{
			Enabled: true,
			TopK:    3,
		},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Content != "answer" {
		t.Fatalf("expected answer, got %s", resp.Content)
	}
	if len(resp.Sources) != 1 {
		t.Fatalf("expected 1 source, got %d", len(resp.Sources))
	}
	if resp.Model != "mock-model" {
		t.Fatalf("expected mock-model, got %s", resp.Model)
	}
}

func TestQueryOrchestratorHandleWithoutKnowledgeProvider(t *testing.T) {
	llmProvider := &mockllm.Provider{
		ChatResponse: llm.ChatResponse{
			Content: "fallback answer",
		},
	}

	orchestrator := NewQueryOrchestrator(llmProvider, nil)
	resp, err := orchestrator.Handle(context.Background(), AIRequest{
		TaskType: TaskTypeQA,
		Query:    "hello",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Content != "fallback answer" {
		t.Fatalf("expected fallback answer, got %s", resp.Content)
	}
	if len(resp.Sources) != 0 {
		t.Fatalf("expected no sources, got %d", len(resp.Sources))
	}
}
