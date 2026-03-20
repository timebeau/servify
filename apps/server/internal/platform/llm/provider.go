package llm

import "context"

// ChatChunk is the streaming unit returned by ChatStream.
type ChatChunk struct {
	ContentDelta string `json:"content_delta"`
	ToolCall     *ToolCall `json:"tool_call,omitempty"`
	Done         bool   `json:"done,omitempty"`
}

// LLMProvider defines the contract for model providers.
type LLMProvider interface {
	Chat(ctx context.Context, req ChatRequest) (ChatResponse, error)
	ChatStream(ctx context.Context, req ChatRequest) (<-chan ChatChunk, error)
	Embed(ctx context.Context, texts []string) ([][]float32, error)
	HealthCheck(ctx context.Context) error
}
