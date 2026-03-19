package application

import (
	"strings"

	"servify/apps/server/internal/platform/knowledgeprovider"
	"servify/apps/server/internal/platform/llm"
)

// PromptBuilder assembles vendor-neutral chat messages for the orchestrator.
type PromptBuilder struct{}

func NewPromptBuilder() *PromptBuilder {
	return &PromptBuilder{}
}

func (b *PromptBuilder) Build(req AIRequest, hits []knowledgeprovider.KnowledgeHit) []llm.ChatMessage {
	messages := append([]llm.ChatMessage(nil), req.Messages...)

	if len(req.Messages) > 0 {
		messages = append([]llm.ChatMessage{{Role: "system", Content: b.buildContextPrompt(req.Messages)}}, messages...)
	}
	if len(hits) > 0 {
		messages = append([]llm.ChatMessage{{Role: "system", Content: b.buildKnowledgePrompt(hits)}}, messages...)
	}
	if req.SystemPrompt != "" {
		messages = append([]llm.ChatMessage{{Role: "system", Content: b.buildSystemPrompt(req.SystemPrompt)}}, messages...)
	}
	if req.Query != "" && len(req.Messages) == 0 {
		messages = append(messages, llm.ChatMessage{Role: "user", Content: req.Query})
	}

	return messages
}

func (b *PromptBuilder) buildSystemPrompt(systemPrompt string) string {
	return strings.TrimSpace(systemPrompt)
}

func (b *PromptBuilder) buildContextPrompt(messages []llm.ChatMessage) string {
	var sb strings.Builder
	sb.WriteString("Conversation context:\n")
	for _, message := range messages {
		role := strings.TrimSpace(message.Role)
		if role == "" {
			role = "unknown"
		}
		sb.WriteString("- ")
		sb.WriteString(role)
		sb.WriteString(": ")
		sb.WriteString(strings.TrimSpace(message.Content))
		sb.WriteString("\n")
	}
	return sb.String()
}

func (b *PromptBuilder) buildKnowledgePrompt(hits []knowledgeprovider.KnowledgeHit) string {
	var sb strings.Builder
	sb.WriteString("Knowledge context:\n")
	for _, hit := range hits {
		sb.WriteString("- ")
		sb.WriteString(hit.Title)
		sb.WriteString(": ")
		sb.WriteString(hit.Content)
		sb.WriteString("\n")
	}
	return sb.String()
}
