package application

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"servify/apps/server/internal/models"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateCustomer(ctx context.Context, cmd CreateCustomerCommand) (*models.User, error) {
	if strings.TrimSpace(cmd.Username) == "" {
		return nil, fmt.Errorf("username required")
	}
	cmd.Tags = normalizeTags(cmd.Tags)
	if cmd.Source == "" {
		cmd.Source = "web"
	}
	if cmd.Priority == "" {
		cmd.Priority = "normal"
	}
	return s.repo.CreateCustomer(ctx, cmd)
}

func (s *Service) GetCustomerByID(ctx context.Context, customerID uint) (*models.User, error) {
	return s.repo.GetCustomerByID(ctx, customerID)
}

func (s *Service) UpdateCustomer(ctx context.Context, customerID uint, cmd UpdateCustomerCommand) (*models.User, error) {
	if cmd.Tags != nil {
		tags := normalizeTags(*cmd.Tags)
		cmd.Tags = &tags
	}
	return s.repo.UpdateCustomer(ctx, customerID, cmd)
}

func (s *Service) ListCustomers(ctx context.Context, query ListCustomersQuery) ([]CustomerInfoDTO, int64, error) {
	query.Page = normalizePage(query.Page)
	query.PageSize = normalizePageSize(query.PageSize)
	query.SortBy = normalizeSortBy(query.SortBy)
	query.SortOrder = normalizeSortOrder(query.SortOrder)
	query.Tags = normalizeTags(query.Tags)
	return s.repo.ListCustomers(ctx, query)
}

func (s *Service) GetCustomerActivity(ctx context.Context, customerID uint, limit int) (*CustomerActivityDTO, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.repo.GetCustomerActivity(ctx, customerID, limit)
}

func (s *Service) AddNote(ctx context.Context, customerID uint, content string, authorID uint) error {
	content = strings.TrimSpace(content)
	if content == "" {
		return fmt.Errorf("note required")
	}
	return s.repo.AddNote(ctx, customerID, CustomerNoteDTO{
		AuthorID:  authorID,
		Content:   content,
		CreatedAt: time.Now(),
	})
}

func (s *Service) UpdateTags(ctx context.Context, customerID uint, tags []string) error {
	return s.repo.UpdateTags(ctx, customerID, normalizeTags(tags))
}

func (s *Service) GetStats(ctx context.Context) (*CustomerStatsDTO, error) {
	return s.repo.GetStats(ctx)
}

func (s *Service) RevokeCustomerTokens(ctx context.Context, customerID uint, revokeAt time.Time) (int, error) {
	if customerID == 0 {
		return 0, fmt.Errorf("customer_id required")
	}
	if revokeAt.IsZero() {
		revokeAt = time.Now().UTC()
	}
	return s.repo.RevokeCustomerTokens(ctx, customerID, revokeAt)
}

func normalizeTags(tags []string) []string {
	out := make([]string, 0, len(tags))
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			continue
		}
		if !slices.Contains(out, tag) {
			out = append(out, tag)
		}
	}
	return out
}

func normalizePage(page int) int {
	if page <= 0 {
		return 1
	}
	return page
}

func normalizePageSize(pageSize int) int {
	if pageSize <= 0 {
		return 20
	}
	if pageSize > 200 {
		return 200
	}
	return pageSize
}

func normalizeSortBy(v string) string {
	switch v {
	case "created_at", "updated_at", "name", "email", "status":
		return v
	default:
		return "created_at"
	}
}

func normalizeSortOrder(v string) string {
	if strings.EqualFold(v, "asc") {
		return "asc"
	}
	return "desc"
}
