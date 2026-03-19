package services

import (
	"context"
	"strings"

	"servify/apps/server/internal/models"
	aimodule "servify/apps/server/internal/modules/ai/application"
	"servify/apps/server/internal/platform/knowledgeprovider"
	"servify/apps/server/internal/platform/llm"
)

// OrchestratedAIService adapts the new AI module to the legacy AIServiceInterface.
type OrchestratedAIService struct {
	orchestrator      *aimodule.QueryOrchestrator
	llmProvider       llm.LLMProvider
	knowledgeProvider knowledgeprovider.KnowledgeProvider
}

func NewOrchestratedAIService(llmProvider llm.LLMProvider, knowledgeProvider knowledgeprovider.KnowledgeProvider) *OrchestratedAIService {
	return &OrchestratedAIService{
		orchestrator:      aimodule.NewQueryOrchestrator(llmProvider, knowledgeProvider),
		llmProvider:       llmProvider,
		knowledgeProvider: knowledgeProvider,
	}
}

func (s *OrchestratedAIService) ProcessQuery(ctx context.Context, query string, sessionID string) (*AIResponse, error) {
	resp, err := s.orchestrator.Handle(ctx, aimodule.AIRequest{
		TaskType:       aimodule.TaskTypeQA,
		ConversationID: sessionID,
		Query:          query,
		SystemPrompt:   "你是 Servify 智能客服助手，请基于上下文给出准确、简洁、专业的中文回答。",
		RetrievalPolicy: aimodule.RetrievalPolicy{
			Enabled:   true,
			TopK:      5,
			Threshold: 0.7,
			Strategy:  "hybrid",
		},
	})
	if err != nil {
		return nil, err
	}

	confidence := 0.6
	if len(resp.Sources) > 0 {
		confidence = 0.85
	}

	return &AIResponse{
		Content:    resp.Content,
		Confidence: confidence,
		Source:     "ai",
	}, nil
}

func (s *OrchestratedAIService) ShouldTransferToHuman(query string, sessionHistory []models.Message) bool {
	query = strings.ToLower(query)
	return strings.Contains(query, "人工") ||
		strings.Contains(query, "客服") ||
		strings.Contains(query, "转接")
}

func (s *OrchestratedAIService) GetSessionSummary(messages []models.Message) (string, error) {
	if len(messages) == 0 {
		return "空会话", nil
	}
	last := messages[len(messages)-1]
	if len(last.Content) > 80 {
		return last.Content[:80], nil
	}
	return last.Content, nil
}

func (s *OrchestratedAIService) InitializeKnowledgeBase() {
	// No-op. Knowledge initialization is owned by provider implementations.
}

func (s *OrchestratedAIService) GetStatus(ctx context.Context) map[string]interface{} {
	return map[string]interface{}{
		"type":               "orchestrated",
		"llm_provider":       s.llmProvider != nil,
		"knowledge_provider": s.knowledgeProvider != nil,
	}
}

var _ AIServiceInterface = (*OrchestratedAIService)(nil)
