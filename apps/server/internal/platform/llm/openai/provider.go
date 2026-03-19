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

	"servify/apps/server/internal/platform/llm"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type request struct {
	Model       string           `json:"model"`
	Messages    []requestMessage `json:"messages"`
	Temperature float64          `json:"temperature,omitempty"`
	MaxTokens   int              `json:"max_tokens,omitempty"`
	Stream      bool             `json:"stream,omitempty"`
}

type requestMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type response struct {
	Choices []struct {
		FinishReason string `json:"finish_reason"`
		Message      struct {
			Content string `json:"content"`
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
	payload := request{
		Model:       req.Model,
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
		Stream:      false,
		Messages:    make([]requestMessage, 0, len(req.Messages)),
	}
	if payload.Model == "" {
		payload.Model = "gpt-3.5-turbo"
	}
	for _, msg := range req.Messages {
		payload.Messages = append(payload.Messages, requestMessage{
			Role:    msg.Role,
			Content: msg.Content,
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
		return llm.ChatResponse{}, fmt.Errorf("send openai request: %w", err)
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return llm.ChatResponse{}, fmt.Errorf("read openai response: %w", err)
	}

	if httpResp.StatusCode >= 400 {
		return llm.ChatResponse{}, fmt.Errorf("openai api error [%d]: %s", httpResp.StatusCode, string(respBody))
	}

	var resp response
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return llm.ChatResponse{}, fmt.Errorf("decode openai response: %w", err)
	}
	if resp.Error != nil {
		return llm.ChatResponse{}, fmt.Errorf("openai api error: %s", resp.Error.Message)
	}
	if len(resp.Choices) == 0 {
		return llm.ChatResponse{}, fmt.Errorf("openai api error: empty choices")
	}

	out := llm.ChatResponse{
		Content:      resp.Choices[0].Message.Content,
		Model:        payload.Model,
		FinishReason: resp.Choices[0].FinishReason,
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

func (p *Provider) ChatStream(ctx context.Context, req llm.ChatRequest) (<-chan llm.ChatChunk, error) {
	return nil, fmt.Errorf("openai stream not implemented yet")
}

func (p *Provider) Embed(ctx context.Context, texts []string) ([][]float32, error) {
	return nil, fmt.Errorf("openai embed not implemented yet")
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
		return fmt.Errorf("openai health check failed [%d]: %s", resp.StatusCode, string(body))
	}
	return nil
}
