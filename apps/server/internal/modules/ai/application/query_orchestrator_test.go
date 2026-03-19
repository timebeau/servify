package application

import (
	"context"
	"testing"

	"servify/apps/server/internal/platform/knowledgeprovider"
	mockkp "servify/apps/server/internal/platform/knowledgeprovider/mock"
	"servify/apps/server/internal/platform/llm"
	mockllm "servify/apps/server/internal/platform/llm/mock"
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

func TestRetrieverDisabledReturnsNoHits(t *testing.T) {
	retriever := NewRetriever(&mockkp.Provider{
		Hits: []knowledgeprovider.KnowledgeHit{
			{DocumentID: "doc-1", Title: "Ignored", Content: "Ignored"},
		},
	})

	hits, err := retriever.Retrieve(context.Background(), AIRequest{
		Query: "ignored",
		RetrievalPolicy: RetrievalPolicy{
			Enabled: false,
		},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(hits) != 0 {
		t.Fatalf("expected no hits, got %d", len(hits))
	}
}

func TestQueryOrchestratorGuardrailsRejectBlockedInput(t *testing.T) {
	llmProvider := &mockllm.Provider{
		ChatResponse: llm.ChatResponse{Content: "should not be returned"},
	}
	orchestrator := NewQueryOrchestrator(llmProvider, nil)

	_, err := orchestrator.Handle(context.Background(), AIRequest{
		Query: "please run DROP TABLE users",
	})
	if err == nil {
		t.Fatal("expected guardrail error")
	}

	metrics := orchestrator.Metrics()
	if metrics.QueryCount != 1 || metrics.FallbackCount != 1 {
		t.Fatalf("unexpected metrics: %+v", metrics)
	}
}

func TestQueryOrchestratorSanitizesLongOutputAndTracksMetrics(t *testing.T) {
	longContent := ""
	for i := 0; i < 2500; i++ {
		longContent += "a"
	}
	llmProvider := &mockllm.Provider{
		ChatResponse: llm.ChatResponse{
			Content: longContent,
			Model:   "mock-model",
			TokenUsage: &llm.TokenUsage{
				TotalTokens: 123,
			},
		},
	}
	orchestrator := NewQueryOrchestrator(llmProvider, nil)

	resp, err := orchestrator.Handle(context.Background(), AIRequest{
		Query: "hello",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !resp.Truncated {
		t.Fatal("expected output to be truncated")
	}
	if len(resp.Content) != 2000 {
		t.Fatalf("expected truncated length 2000, got %d", len(resp.Content))
	}

	metrics := orchestrator.Metrics()
	if metrics.QueryCount != 1 || metrics.SuccessCount != 1 {
		t.Fatalf("unexpected metrics: %+v", metrics)
	}
	if metrics.LastTokenUsage != 123 {
		t.Fatalf("unexpected token usage: %+v", metrics)
	}
}
