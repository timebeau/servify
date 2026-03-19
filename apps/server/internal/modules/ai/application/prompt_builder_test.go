package application

import (
	"strings"
	"testing"

	"servify/apps/server/internal/platform/knowledgeprovider"
	"servify/apps/server/internal/platform/llm"
)

func TestPromptBuilderBuildIncludesContextPrompt(t *testing.T) {
	builder := NewPromptBuilder()
	messages := builder.Build(AIRequest{
		SystemPrompt: "system",
		Messages: []llm.ChatMessage{
			{Role: "user", Content: "hello"},
			{Role: "assistant", Content: "hi"},
		},
	}, []knowledgeprovider.KnowledgeHit{{Title: "KB", Content: "fact"}})

	if len(messages) < 4 {
		t.Fatalf("expected prompt messages, got %d", len(messages))
	}
	found := false
	for _, message := range messages {
		if strings.Contains(message.Content, "Conversation context:") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected context prompt, got %+v", messages)
	}
}
