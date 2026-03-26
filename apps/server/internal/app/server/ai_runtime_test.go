package server

import (
	"testing"

	"servify/apps/server/internal/config"
)

func TestAIAssemblyKnowledgeProviderUsesResolvedKnowledgeBaseID(t *testing.T) {
	assembly := &AIAssembly{
		WeKnoraHealthy:  true,
		WeKnoraClient:   nil,
		KnowledgeBaseID: "kb-resolved",
	}
	if provider := assembly.KnowledgeProvider(&config.Config{}); provider != nil {
		t.Fatal("expected nil provider without client")
	}
	if assembly.KnowledgeBaseID != "kb-resolved" {
		t.Fatalf("knowledge base id = %q want kb-resolved", assembly.KnowledgeBaseID)
	}
}
