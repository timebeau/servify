package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"servify/apps/server/internal/config"
	"servify/apps/server/internal/platform/llm"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type request struct {
	Model       string           `json:"model"`
	Messages    []requestMessage `json:"messages"`
	Tools       []requestTool    `json:"tools,omitempty"`
	Temperature float64          `json:"temperature,omitempty"`
	MaxTokens   int              `json:"max_tokens,omitempty"`
	Stream      bool             `json:"stream,omitempty"`
}

type requestMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type requestTool struct {
	Type     string              `json:"type"`
	Function requestToolFunction `json:"function"`
}

type requestToolFunction struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
}

type response struct {
	Choices []struct {
		FinishReason string `json:"finish_reason"`
		Message      struct {
			Content   string `json:"content"`
			ToolCalls []struct {
				ID       string `json:"id"`
				Type     string `json:"type"`
				Function struct {
					Name      string `json:"name"`
					Arguments string `json:"arguments"`
				} `json:"function"`
			} `json:"tool_calls"`
		} `json:"message"`
	} `json:"choices"`
	Usage *struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

// Provider is the OpenAI implementation of llm.LLMProvider.
type Provider struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func NewProvider(apiKey, baseURL string) *Provider {
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	return &Provider{
		apiKey:  apiKey,
		baseURL: strings.TrimRight(baseURL, "/"),
		client: &http.Client{
			Timeout:   30 * time.Second,
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		},
	}
}

func (p *Provider) Chat(ctx context.Context, req llm.ChatRequest) (llm.ChatResponse, error) {
	var out llm.ChatResponse
	callCtx, cancel := llm.WithRequestTimeout(ctx, req.Options, 30*time.Second)
	defer cancel()

	err := llm.Retry(callCtx, req.Options.RetryPolicy, func(runCtx context.Context) error {
		resp, err := p.chatOnce(runCtx, req)
		if err != nil {
			return err
		}
		out = resp
		return nil
	})
	return out, err
}

func (p *Provider) ChatStream(ctx context.Context, req llm.ChatRequest) (<-chan llm.ChatChunk, error) {
	return nil, &llm.ProviderError{
		Provider:  "openai",
		Code:      llm.ProviderErrorNotSupported,
		Message:   "openai stream not implemented yet",
		Retryable: false,
	}
}

func (p *Provider) Embed(ctx context.Context, texts []string) ([][]float32, error) {
	return nil, &llm.ProviderError{
		Provider:  "openai",
		Code:      llm.ProviderErrorNotSupported,
		Message:   "openai embed not implemented yet",
		Retryable: false,
	}
}

func (p *Provider) HealthCheck(ctx context.Context) error {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, p.baseURL+"/models", nil)
	if err != nil {
		return err
	}
	if p.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)
	}
	resp, err := p.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return llm.HTTPError("openai", resp.StatusCode, string(body))
	}
	return nil
}

func (p *Provider) chatOnce(ctx context.Context, req llm.ChatRequest) (llm.ChatResponse, error) {
	payload := request{
		Model:       req.Model,
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
		Stream:      false,
		Messages:    make([]requestMessage, 0, len(req.Messages)),
		Tools:       make([]requestTool, 0, len(req.Tools)),
	}
	if payload.Model == "" {
		payload.Model = config.DefaultOpenAIModel
	}
	for _, msg := range req.Messages {
		payload.Messages = append(payload.Messages, requestMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}
	for _, tool := range req.Tools {
		payload.Tools = append(payload.Tools, requestTool{
			Type: "function",
			Function: requestToolFunction{
				Name:        tool.Name,
				Description: tool.Description,
				Parameters:  tool.InputSchema,
			},
		})
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return llm.ChatResponse{}, fmt.Errorf("marshal openai request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/chat/completions", bytes.NewBuffer(body))
	if err != nil {
		return llm.ChatResponse{}, fmt.Errorf("create openai request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if p.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)
	}

	httpResp, err := p.client.Do(httpReq)
	if err != nil {
		return llm.ChatResponse{}, &llm.ProviderError{
			Provider:  "openai",
			Code:      llm.ProviderErrorUnavailable,
			Message:   "send openai request failed",
			Retryable: true,
			Cause:     err,
		}
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return llm.ChatResponse{}, fmt.Errorf("read openai response: %w", err)
	}

	if httpResp.StatusCode >= 400 {
		return llm.ChatResponse{}, llm.HTTPError("openai", httpResp.StatusCode, string(respBody))
	}

	var resp response
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return llm.ChatResponse{}, fmt.Errorf("decode openai response: %w", err)
	}
	if resp.Error != nil {
		return llm.ChatResponse{}, &llm.ProviderError{
			Provider:  "openai",
			Code:      llm.ProviderErrorUpstream,
			Message:   resp.Error.Message,
			Retryable: false,
		}
	}
	if len(resp.Choices) == 0 {
		return llm.ChatResponse{}, &llm.ProviderError{
			Provider:  "openai",
			Code:      llm.ProviderErrorUpstream,
			Message:   "empty choices",
			Retryable: false,
		}
	}

	out := llm.ChatResponse{
		Content:      resp.Choices[0].Message.Content,
		Provider:     "openai",
		Model:        payload.Model,
		FinishReason: resp.Choices[0].FinishReason,
		ToolCalls:    decodeToolCalls(resp.Choices[0].Message.ToolCalls),
	}
	if resp.Usage != nil {
		out.TokenUsage = &llm.TokenUsage{
			InputTokens:  resp.Usage.PromptTokens,
			OutputTokens: resp.Usage.CompletionTokens,
			TotalTokens:  resp.Usage.TotalTokens,
		}
	}
	return out, nil
}

func decodeToolCalls(toolCalls []struct {
	ID       string "json:\"id\""
	Type     string "json:\"type\""
	Function struct {
		Name      string "json:\"name\""
		Arguments string "json:\"arguments\""
	} "json:\"function\""
}) []llm.ToolCall {
	if len(toolCalls) == 0 {
		return nil
	}
	out := make([]llm.ToolCall, 0, len(toolCalls))
	for _, toolCall := range toolCalls {
		args := map[string]interface{}{}
		if strings.TrimSpace(toolCall.Function.Arguments) != "" {
			_ = json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
		}
		out = append(out, llm.ToolCall{
			ID:        toolCall.ID,
			Name:      toolCall.Function.Name,
			Arguments: args,
		})
	}
	return out
}
