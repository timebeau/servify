package mock

import (
	"context"

	"servify/apps/server/internal/platform/llm"
)

// Provider is a controllable mock implementation of llm.LLMProvider.
type Provider struct {
	ChatResponse      llm.ChatResponse
	ChatError         error
	StreamChunks      []llm.ChatChunk
	StreamError       error
	EmbeddingResponse [][]float32
	EmbeddingError    error
	HealthError       error
}

func (p *Provider) Chat(ctx context.Context, req llm.ChatRequest) (llm.ChatResponse, error) {
	return p.ChatResponse, p.ChatError
}

func (p *Provider) ChatStream(ctx context.Context, req llm.ChatRequest) (<-chan llm.ChatChunk, error) {
	if p.StreamError != nil {
		return nil, p.StreamError
	}
	ch := make(chan llm.ChatChunk, len(p.StreamChunks))
	for _, chunk := range p.StreamChunks {
		ch <- chunk
	}
	close(ch)
	return ch, nil
}

func (p *Provider) Embed(ctx context.Context, texts []string) ([][]float32, error) {
	return p.EmbeddingResponse, p.EmbeddingError
}

func (p *Provider) HealthCheck(ctx context.Context) error {
	return p.HealthError
}
