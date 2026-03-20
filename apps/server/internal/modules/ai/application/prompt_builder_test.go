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

func TestPromptBuilderBuildAppendsQueryWhenNoMessages(t *testing.T) {
	builder := NewPromptBuilder()
	messages := builder.Build(AIRequest{
		SystemPrompt: "  system  ",
		Query:        "hello world",
	}, nil)

	if len(messages) != 2 {
		t.Fatalf("expected system and user prompt, got %d", len(messages))
	}
	if messages[0].Role != "system" || messages[0].Content != "system" {
		t.Fatalf("unexpected system prompt: %+v", messages[0])
	}
	if messages[1].Role != "user" || messages[1].Content != "hello world" {
		t.Fatalf("unexpected user prompt: %+v", messages[1])
	}
}

func TestPromptBuilderBuildOrdersKnowledgeAndContextBeforeMessages(t *testing.T) {
	builder := NewPromptBuilder()
	messages := builder.Build(AIRequest{
		SystemPrompt: "system",
		Messages: []llm.ChatMessage{
			{Role: "", Content: "user question"},
		},
	}, []knowledgeprovider.KnowledgeHit{
		{Title: "KB", Content: "fact"},
	})

	if len(messages) < 4 {
		t.Fatalf("expected at least 4 messages, got %d", len(messages))
	}
	if messages[0].Content != "system" {
		t.Fatalf("expected system prompt first, got %+v", messages[0])
	}
	if !strings.Contains(messages[1].Content, "Knowledge context:") {
		t.Fatalf("expected knowledge context second, got %+v", messages[1])
	}
	if !strings.Contains(messages[2].Content, "Conversation context:") {
		t.Fatalf("expected conversation context third, got %+v", messages[2])
	}
}
