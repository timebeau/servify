package delivery

import (
	"context"

	suggestionapp "servify/apps/server/internal/modules/suggestion/application"
	suggestioncontract "servify/apps/server/internal/modules/suggestion/contract"
	suggestioninfra "servify/apps/server/internal/modules/suggestion/infra"

	"gorm.io/gorm"
)

type HandlerServiceAdapter struct {
	service *suggestionapp.Service
}

func NewHandlerService(db *gorm.DB) *HandlerServiceAdapter {
	return NewHandlerServiceAdapter(suggestionapp.NewService(suggestioninfra.NewGormRepository(db)))
}

func NewHandlerServiceAdapter(service *suggestionapp.Service) *HandlerServiceAdapter {
	return &HandlerServiceAdapter{service: service}
}

func (a *HandlerServiceAdapter) Suggest(ctx context.Context, req *suggestioncontract.SuggestionRequest) (*suggestioncontract.SuggestionResponse, error) {
	return a.service.Suggest(ctx, req)
}

var _ HandlerService = (*HandlerServiceAdapter)(nil)
