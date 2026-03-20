package delivery

import (
	"context"

	"servify/apps/server/internal/models"
	"servify/apps/server/internal/services"
)

// HandlerServiceAdapter bridges the legacy knowledge service to the handler-facing contract.
type HandlerServiceAdapter struct {
	service *services.KnowledgeDocService
}

func NewHandlerServiceAdapter(service *services.KnowledgeDocService) *HandlerServiceAdapter {
	return &HandlerServiceAdapter{service: service}
}

func (a *HandlerServiceAdapter) List(ctx context.Context, req *services.KnowledgeDocListRequest) ([]models.KnowledgeDoc, int64, error) {
	return a.service.List(ctx, req)
}

func (a *HandlerServiceAdapter) Get(ctx context.Context, id uint) (*models.KnowledgeDoc, error) {
	return a.service.Get(ctx, id)
}

func (a *HandlerServiceAdapter) Create(ctx context.Context, req *services.KnowledgeDocCreateRequest) (*models.KnowledgeDoc, error) {
	return a.service.Create(ctx, req)
}

func (a *HandlerServiceAdapter) Update(ctx context.Context, id uint, req *services.KnowledgeDocUpdateRequest) (*models.KnowledgeDoc, error) {
	return a.service.Update(ctx, id, req)
}

func (a *HandlerServiceAdapter) Delete(ctx context.Context, id uint) error {
	return a.service.Delete(ctx, id)
}
