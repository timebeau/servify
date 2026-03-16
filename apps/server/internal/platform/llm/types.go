package llm

// ChatMessage is the vendor-neutral message format used by LLM providers.
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ToolCall describes a tool request returned by a model.
type ToolCall struct {
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

// TokenUsage tracks provider token consumption.
type TokenUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

// ChatRequest is the shared input model for chat-style model requests.
type ChatRequest struct {
	Model       string        `json:"model,omitempty"`
	Messages    []ChatMessage `json:"messages"`
	Temperature float64       `json:"temperature,omitempty"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Stream      bool          `json:"stream,omitempty"`
}

// ChatResponse is the shared output model for chat-style model responses.
type ChatResponse struct {
	Content      string      `json:"content"`
	Model        string      `json:"model,omitempty"`
	FinishReason string      `json:"finish_reason,omitempty"`
	ToolCalls    []ToolCall  `json:"tool_calls,omitempty"`
	TokenUsage   *TokenUsage `json:"token_usage,omitempty"`
}
