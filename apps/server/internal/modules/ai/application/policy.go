package application

import "context"

type PolicyDecision struct {
	Allowed bool
	Reason  string
}

type PolicyHook interface {
	Evaluate(ctx context.Context, req AIRequest) (PolicyDecision, error)
}

type PromptAuditRecord struct {
	PromptVersion   string
	TaskType        TaskType
	Provider        string
	MessageCount    int
	RetrievalHits   int
	SystemPromptSet bool
}

type PromptAuditRecorder interface {
	RecordPrompt(ctx context.Context, record PromptAuditRecord) error
}
