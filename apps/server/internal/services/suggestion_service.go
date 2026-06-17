package services

import (
	"context"

	suggestionapp "servify/apps/server/internal/modules/suggestion/application"
	suggestioncontract "servify/apps/server/internal/modules/suggestion/contract"
	suggestiondelivery "servify/apps/server/internal/modules/suggestion/delivery"

	"gorm.io/gorm"
)

type SuggestionService struct {
	service suggestiondelivery.HandlerService
}

func NewSuggestionService(db *gorm.DB) *SuggestionService {
	return &SuggestionService{service: suggestiondelivery.NewHandlerService(db)}
}

type IntentSuggestion = suggestioncontract.IntentSuggestion
type TicketSuggestion = suggestioncontract.TicketSuggestion
type KnowledgeDocSuggestion = suggestioncontract.KnowledgeDocSuggestion
type SuggestionResponse = suggestioncontract.SuggestionResponse
type SuggestionRequest = suggestioncontract.SuggestionRequest

func (s *SuggestionService) Suggest(ctx context.Context, req *SuggestionRequest) (*SuggestionResponse, error) {
	return s.service.Suggest(ctx, req)
}

var (
	extractTokens        = suggestionapp.ExtractTokens
	buildLikeWhereTokens = suggestionapp.BuildLikeWhereTokens
	scoreText            = suggestionapp.ScoreText
	classifyIntent       = suggestionapp.ClassifyIntent
)
