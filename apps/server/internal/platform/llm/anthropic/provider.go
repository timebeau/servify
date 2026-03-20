package anthropic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"servify/apps/server/internal/platform/llm"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type request struct {
	Model       string    `json:"model"`
	Messages    []message `json:"messages"`
	System      string    `json:"system,omitempty"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature,omitempty"`
	Tools       []tool    `json:"tools,omitempty"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	InputSchema map[string]interface{} `json:"input_schema,omitempty"`
}

type response struct {
	ID         string `json:"id"`
	Model      string `json:"model"`
	StopReason string `json:"stop_reason"`
	Content    []struct {
		Type  string                 `json:"type"`
		Text  string                 `json:"text,omitempty"`
		ID    string                 `json:"id,omitempty"`
		Name  string                 `json:"name,omitempty"`
		Input map[string]interface{} `json:"input,omitempty"`
	} `json:"content"`
	Usage *struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
	Error *struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
}

type Provider struct {
	apiKey  string
	baseURL string
	version string
	client  *http.Client
}

func NewProvider(apiKey, baseURL string) *Provider {
	if baseURL == "" {
		baseURL = "https://api.anthropic.com/v1"
	}
	return &Provider{
		apiKey:  apiKey,
		baseURL: strings.TrimRight(baseURL, "/"),
		version: "2023-06-01",
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
		Provider:  "anthropic",
		Code:      llm.ProviderErrorNotSupported,
		Message:   "anthropic stream not implemented yet",
		Retryable: false,
	}
}

func (p *Provider) Embed(ctx context.Context, texts []string) ([][]float32, error) {
	return nil, &llm.ProviderError{
		Provider:  "anthropic",
		Code:      llm.ProviderErrorNotSupported,
		Message:   "anthropic embed not implemented yet",
		Retryable: false,
	}
}

func (p *Provider) HealthCheck(ctx context.Context) error {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, p.baseURL+"/models", nil)
	if err != nil {
		return err
	}
	if p.apiKey != "" {
		httpReq.Header.Set("x-api-key", p.apiKey)
	}
	httpReq.Header.Set("anthropic-version", p.version)
	resp, err := p.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return llm.HTTPError("anthropic", resp.StatusCode, string(body))
	}
	return nil
}

func (p *Provider) chatOnce(ctx context.Context, req llm.ChatRequest) (llm.ChatResponse, error) {
	payload := request{
		Model:       req.Model,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		Messages:    make([]message, 0, len(req.Messages)),
		Tools:       make([]tool, 0, len(req.Tools)),
	}
	if payload.Model == "" {
		payload.Model = "claude-3-haiku-20240307"
	}
	if payload.MaxTokens == 0 {
		payload.MaxTokens = 1024
	}
	for _, msg := range req.Messages {
		if msg.Role == "system" {
			payload.System = strings.TrimSpace(strings.Join([]string{payload.System, msg.Content}, "\n"))
			continue
		}
		payload.Messages = append(payload.Messages, message{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}
	for _, t := range req.Tools {
		payload.Tools = append(payload.Tools, tool{
			Name:        t.Name,
			Description: t.Description,
			InputSchema: t.InputSchema,
		})
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return llm.ChatResponse{}, fmt.Errorf("marshal anthropic request: %w", err)
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/messages", bytes.NewBuffer(body))
	if err != nil {
		return llm.ChatResponse{}, fmt.Errorf("create anthropic request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("anthropic-version", p.version)
	if p.apiKey != "" {
		httpReq.Header.Set("x-api-key", p.apiKey)
	}

	httpResp, err := p.client.Do(httpReq)
	if err != nil {
		return llm.ChatResponse{}, &llm.ProviderError{
			Provider:  "anthropic",
			Code:      llm.ProviderErrorUnavailable,
			Message:   "send anthropic request failed",
			Retryable: true,
			Cause:     err,
		}
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return llm.ChatResponse{}, fmt.Errorf("read anthropic response: %w", err)
	}
	if httpResp.StatusCode >= 400 {
		return llm.ChatResponse{}, llm.HTTPError("anthropic", httpResp.StatusCode, string(respBody))
	}

	var resp response
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return llm.ChatResponse{}, fmt.Errorf("decode anthropic response: %w", err)
	}
	if resp.Error != nil {
		return llm.ChatResponse{}, &llm.ProviderError{
			Provider:  "anthropic",
			Code:      llm.ProviderErrorUpstream,
			Message:   resp.Error.Message,
			Retryable: false,
		}
	}

	out := llm.ChatResponse{
		Provider:     "anthropic",
		Model:        resp.Model,
		FinishReason: resp.StopReason,
	}
	for _, item := range resp.Content {
		switch item.Type {
		case "text":
			out.Content += item.Text
		case "tool_use":
			out.ToolCalls = append(out.ToolCalls, llm.ToolCall{
				ID:        item.ID,
				Name:      item.Name,
				Arguments: item.Input,
			})
		}
	}
	if resp.Usage != nil {
		out.TokenUsage = &llm.TokenUsage{
			InputTokens:  resp.Usage.InputTokens,
			OutputTokens: resp.Usage.OutputTokens,
			TotalTokens:  resp.Usage.InputTokens + resp.Usage.OutputTokens,
		}
	}
	return out, nil
}
