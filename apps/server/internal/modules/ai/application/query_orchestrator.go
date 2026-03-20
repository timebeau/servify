package application

import (
	"context"
	"time"

	"servify/apps/server/internal/platform/knowledgeprovider"
	"servify/apps/server/internal/platform/llm"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// QueryOrchestrator coordinates retrieval and model execution.
type QueryOrchestrator struct {
	llmProvider   llm.LLMProvider
	retriever     *Retriever
	promptBuilder *PromptBuilder
	guardrails    *Guardrails
	metrics       *Metrics
	tracer        trace.Tracer
	policyHooks   []PolicyHook
	auditRecorder PromptAuditRecorder
}

func NewQueryOrchestrator(llmProvider llm.LLMProvider, knowledgeProvider knowledgeprovider.KnowledgeProvider) *QueryOrchestrator {
	return &QueryOrchestrator{
		llmProvider:   llmProvider,
		retriever:     NewRetriever(knowledgeProvider),
		promptBuilder: NewPromptBuilder(),
		guardrails:    NewGuardrails(),
		metrics:       NewMetrics(),
		tracer:        otel.Tracer("servify.ai.orchestrator"),
	}
}

func (o *QueryOrchestrator) Handle(ctx context.Context, req AIRequest) (*AIResponse, error) {
	ctx, span := o.tracer.Start(ctx, "ai.query_orchestrator.handle")
	defer span.End()

	span.SetAttributes(
		attribute.String("ai.task_type", string(req.TaskType)),
		attribute.Bool("ai.retrieval.enabled", req.RetrievalPolicy.Enabled),
	)

	start := time.Now()
	o.metrics.RecordQuery()
	if o.llmProvider == nil {
		return nil, nil
	}
	for _, hook := range o.policyHooks {
		decision, err := hook.Evaluate(ctx, req)
		if err != nil {
			o.metrics.RecordError("policy", "policy_hook_error")
			o.metrics.RecordFallback()
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}
		if !decision.Allowed {
			o.metrics.RecordPolicyRejection(decision.Reason)
			o.metrics.RecordFallback()
			span.SetStatus(codes.Error, decision.Reason)
			return nil, context.Canceled
		}
	}
	if err := o.guardrails.ValidateInput(req); err != nil {
		o.metrics.RecordPolicyRejection("guardrails")
		o.metrics.RecordFallback()
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	hits, err := o.retriever.Retrieve(ctx, req)
	if err != nil {
		o.metrics.RecordError("knowledge", "retrieval")
		o.metrics.RecordFallback()
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	messages := o.promptBuilder.Build(req, hits)
	if o.auditRecorder != nil {
		_ = o.auditRecorder.RecordPrompt(ctx, PromptAuditRecord{
			PromptVersion:   "v1",
			TaskType:        req.TaskType,
			MessageCount:    len(messages),
			RetrievalHits:   len(hits),
			SystemPromptSet: req.SystemPrompt != "",
		})
	}

	chatResp, err := o.llmProvider.Chat(ctx, llm.ChatRequest{
		Messages: messages,
	})
	if err != nil {
		o.metrics.RecordError("llm", "chat")
		o.metrics.RecordFallback()
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	content, truncated := o.guardrails.SanitizeOutput(chatResp.Content)
	totalTokens := 0
	if chatResp.TokenUsage != nil {
		totalTokens = chatResp.TokenUsage.TotalTokens
	}
	provider := "llm"
	if chatResp.Provider != "" {
		provider = chatResp.Provider
	}
	o.metrics.RecordSuccess(provider, time.Since(start), totalTokens)
	span.SetAttributes(
		attribute.String("ai.provider", provider),
		attribute.Int("ai.tokens.total", totalTokens),
		attribute.Int64("ai.latency.ms", time.Since(start).Milliseconds()),
	)

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

func (o *QueryOrchestrator) SetPolicyHooks(hooks ...PolicyHook) {
	o.policyHooks = append([]PolicyHook(nil), hooks...)
}

func (o *QueryOrchestrator) SetAuditRecorder(recorder PromptAuditRecorder) {
	o.auditRecorder = recorder
}
