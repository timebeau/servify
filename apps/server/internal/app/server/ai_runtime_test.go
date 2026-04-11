package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"servify/apps/server/internal/config"

	"github.com/sirupsen/logrus"
)

func TestAIAssemblyKnowledgeProviderUsesResolvedKnowledgeBaseID(t *testing.T) {
	assembly := &AIAssembly{
		KnowledgeProviderHealthy: true,
		WeKnoraHealthy:           true,
		WeKnoraClient:            nil,
		KnowledgeBaseID:          "kb-resolved",
	}
	if provider := assembly.KnowledgeProvider(&config.Config{}); provider != nil {
		t.Fatal("expected nil provider without client")
	}
	if assembly.KnowledgeBaseID != "kb-resolved" {
		t.Fatalf("knowledge base id = %q want kb-resolved", assembly.KnowledgeBaseID)
	}
}

func TestBuildAIAssemblyPrefersDifyOverWeKnora(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/datasets/ds-1":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":"ds-1","name":"Primary Dify Dataset"}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	cfg := config.GetDefaultConfig()
	cfg.Dify.Enabled = true
	cfg.Dify.BaseURL = server.URL
	cfg.Dify.APIKey = "dify-key"
	cfg.Dify.DatasetID = "ds-1"
	cfg.WeKnora.Enabled = true
	cfg.WeKnora.BaseURL = "http://127.0.0.1:1"
	cfg.WeKnora.APIKey = "wk-key"

	assembly, err := BuildAIAssembly(cfg, logrus.New(), AIAssemblyOptions{})
	if err != nil {
		t.Fatalf("BuildAIAssembly() error = %v", err)
	}
	if assembly.KnowledgeProviderID != "dify" {
		t.Fatalf("knowledge provider = %q", assembly.KnowledgeProviderID)
	}
	if !assembly.KnowledgeProviderHealthy {
		t.Fatalf("expected knowledge provider health to be tracked")
	}
	status := assembly.RuntimeService.GetStatus(context.Background())
	if provider, _ := status["knowledge_provider"].(string); provider != "dify" {
		t.Fatalf("status provider = %q", provider)
	}
}
