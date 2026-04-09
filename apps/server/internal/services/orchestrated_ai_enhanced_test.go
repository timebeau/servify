package services

import (
	"context"
	"testing"

	"servify/apps/server/internal/platform/knowledgeprovider"
	mockkp "servify/apps/server/internal/platform/knowledgeprovider/mock"
	"servify/apps/server/internal/platform/llm"
	mockllm "servify/apps/server/internal/platform/llm/mock"
)

func TestOrchestratedEnhancedAIServiceProcessQueryEnhanced(t *testing.T) {
	base := NewAIService("", "")
	base.InitializeKnowledgeBase()

	svc := NewOrchestratedEnhancedAIService(
		base,
		&mockllm.Provider{
			ChatResponse: llm.ChatResponse{Content: "answer", TokenUsage: &llm.TokenUsage{TotalTokens: 42}},
		},
		&mockkp.Provider{
			Hits: []knowledgeprovider.KnowledgeHit{
				{DocumentID: "doc-1", Title: "Billing", Content: "Billing details", Score: 0.91},
			},
		},
		"",
		nil,
		"kb-1",
		nil,
	)

	resp, err := svc.ProcessQueryEnhanced(context.Background(), "billing", "session-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Strategy != "weknora" {
		t.Fatalf("expected weknora strategy, got %s", resp.Strategy)
	}
	if len(resp.Sources) != 1 {
		t.Fatalf("expected sources, got %+v", resp.Sources)
	}
	if resp.TokensUsed != 42 {
		t.Fatalf("expected tokens used to be tracked, got %d", resp.TokensUsed)
	}

	metrics := svc.GetMetrics()
	if metrics.QueryCount != 1 || metrics.SuccessCount != 1 || metrics.WeKnoraUsageCount != 1 {
		t.Fatalf("unexpected metrics: %+v", metrics)
	}
}

func TestOrchestratedEnhancedAIServiceFallbackAndReset(t *testing.T) {
	base := NewAIService("", "")
	base.InitializeKnowledgeBase()

	svc := NewOrchestratedEnhancedAIService(
		base,
		&mockllm.Provider{ChatError: context.DeadlineExceeded},
		&mockkp.Provider{},
		"",
		nil,
		"kb-1",
		nil,
	)

	resp, err := svc.ProcessQueryEnhanced(context.Background(), "Servify 是什么", "session-1")
	if err != nil {
		t.Fatalf("expected fallback response, got %v", err)
	}
	if resp.Strategy != "fallback" {
		t.Fatalf("expected fallback strategy, got %s", resp.Strategy)
	}

	status := svc.GetStatus(context.Background())
	cb := status["circuit_breaker"].(map[string]interface{})
	if cb["failure_count"].(int) == 0 {
		t.Fatalf("expected circuit breaker failures, got %+v", cb)
	}

	metrics := svc.GetMetrics()
	if metrics.FallbackUsageCount == 0 {
		t.Fatalf("expected fallback metrics, got %+v", metrics)
	}

	svc.ResetCircuitBreaker()
	status = svc.GetStatus(context.Background())
	cb = status["circuit_breaker"].(map[string]interface{})
	if cb["failure_count"].(int) != 0 {
		t.Fatalf("expected reset circuit breaker, got %+v", cb)
	}
}
