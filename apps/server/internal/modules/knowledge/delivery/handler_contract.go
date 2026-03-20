package delivery

import (
	"context"

	"servify/apps/server/internal/models"
	"servify/apps/server/internal/services"
)

// HandlerService is the only knowledge contract that HTTP handlers should depend on.
type HandlerService interface {
	List(ctx context.Context, req *services.KnowledgeDocListRequest) ([]models.KnowledgeDoc, int64, error)
	Get(ctx context.Context, id uint) (*models.KnowledgeDoc, error)
	Create(ctx context.Context, req *services.KnowledgeDocCreateRequest) (*models.KnowledgeDoc, error)
	Update(ctx context.Context, id uint, req *services.KnowledgeDocUpdateRequest) (*models.KnowledgeDoc, error)
	Delete(ctx context.Context, id uint) error
}
