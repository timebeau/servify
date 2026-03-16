package application

import (
	"context"
	"strings"
	"time"

	"servify/apps/server/internal/platform/knowledgeprovider"
	"servify/apps/server/internal/platform/llm"
)

// QueryOrchestrator coordinates retrieval and model execution.
type QueryOrchestrator struct {
	llmProvider       llm.LLMProvider
	knowledgeProvider knowledgeprovider.KnowledgeProvider
}

func NewQueryOrchestrator(llmProvider llm.LLMProvider, knowledgeProvider knowledgeprovider.KnowledgeProvider) *QueryOrchestrator {
	return &QueryOrchestrator{
		llmProvider:       llmProvider,
		knowledgeProvider: knowledgeProvider,
	}
}

func (o *QueryOrchestrator) Handle(ctx context.Context, req AIRequest) (*AIResponse, error) {
	start := time.Now()

	messages := append([]llm.ChatMessage(nil), req.Messages...)
	var hits []knowledgeprovider.KnowledgeHit

	if req.RetrievalPolicy.Enabled && o.knowledgeProvider != nil && strings.TrimSpace(req.Query) != "" {
		searchReq := knowledgeprovider.SearchRequest{
			Query:       req.Query,
			TopK:        req.RetrievalPolicy.TopK,
			Threshold:   req.RetrievalPolicy.Threshold,
			Strategy:    req.RetrievalPolicy.Strategy,
			KnowledgeID: req.TenantID,
		}
		found, err := o.knowledgeProvider.Search(ctx, searchReq)
		if err != nil {
			return nil, err
		}
		hits = found
	}

	if req.SystemPrompt != "" {
		messages = append([]llm.ChatMessage{{Role: "system", Content: req.SystemPrompt}}, messages...)
	}
	if req.Query != "" && len(req.Messages) == 0 {
		messages = append(messages, llm.ChatMessage{Role: "user", Content: req.Query})
	}
	if len(hits) > 0 {
		var sb strings.Builder
		sb.WriteString("Knowledge context:\n")
		for _, hit := range hits {
			sb.WriteString("- ")
			sb.WriteString(hit.Title)
			sb.WriteString(": ")
			sb.WriteString(hit.Content)
			sb.WriteString("\n")
		}
		messages = append([]llm.ChatMessage{{Role: "system", Content: sb.String()}}, messages...)
	}

	chatResp, err := o.llmProvider.Chat(ctx, llm.ChatRequest{
		Messages: messages,
	})
	if err != nil {
		return nil, err
	}

	return &AIResponse{
		Content:      chatResp.Content,
		Model:        chatResp.Model,
		Sources:      hits,
		TokenUsage:   chatResp.TokenUsage,
		FinishReason: chatResp.FinishReason,
		Latency:      time.Since(start),
	}, nil
}
