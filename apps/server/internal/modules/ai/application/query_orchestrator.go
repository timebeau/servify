package application

import (
	"context"
	"time"

	"servify/apps/server/internal/platform/knowledgeprovider"
	"servify/apps/server/internal/platform/llm"
)

// QueryOrchestrator coordinates retrieval and model execution.
type QueryOrchestrator struct {
	llmProvider       llm.LLMProvider
	retriever         *Retriever
	promptBuilder     *PromptBuilder
	guardrails        *Guardrails
	metrics           *Metrics
}

func NewQueryOrchestrator(llmProvider llm.LLMProvider, knowledgeProvider knowledgeprovider.KnowledgeProvider) *QueryOrchestrator {
	return &QueryOrchestrator{
		llmProvider:   llmProvider,
		retriever:     NewRetriever(knowledgeProvider),
		promptBuilder: NewPromptBuilder(),
		guardrails:    NewGuardrails(),
		metrics:       NewMetrics(),
	}
}

func (o *QueryOrchestrator) Handle(ctx context.Context, req AIRequest) (*AIResponse, error) {
	start := time.Now()
	o.metrics.RecordQuery()
	if o.llmProvider == nil {
		return nil, nil
	}
	if err := o.guardrails.ValidateInput(req); err != nil {
		o.metrics.RecordFallback()
		return nil, err
	}

	hits, err := o.retriever.Retrieve(ctx, req)
	if err != nil {
		o.metrics.RecordFallback()
		return nil, err
	}
	messages := o.promptBuilder.Build(req, hits)

	chatResp, err := o.llmProvider.Chat(ctx, llm.ChatRequest{
		Messages: messages,
	})
	if err != nil {
		o.metrics.RecordFallback()
		return nil, err
	}
	content, truncated := o.guardrails.SanitizeOutput(chatResp.Content)
	totalTokens := 0
	if chatResp.TokenUsage != nil {
		totalTokens = chatResp.TokenUsage.TotalTokens
	}
	provider := "llm"
	o.metrics.RecordSuccess(provider, time.Since(start), totalTokens)

	return &AIResponse{
		Content:      content,
		Model:        chatResp.Model,
		Provider:     provider,
		Sources:      hits,
		TokenUsage:   chatResp.TokenUsage,
		FinishReason: chatResp.FinishReason,
		Latency:      time.Since(start),
		Truncated:    truncated,
	}, nil
}

func (o *QueryOrchestrator) Metrics() MetricsSnapshot {
	return o.metrics.Snapshot()
}
