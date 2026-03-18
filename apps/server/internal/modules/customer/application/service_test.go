package application

import (
	"context"
	"testing"

	"servify/apps/server/internal/models"
)

type stubRepo struct {
	lastCreate CreateCustomerCommand
	lastUpdate UpdateCustomerCommand
	lastTags   []string
	lastNote   CustomerNoteDTO
}

func (s *stubRepo) CreateCustomer(ctx context.Context, cmd CreateCustomerCommand) (*models.User, error) {
	s.lastCreate = cmd
	return &models.User{ID: 1, Username: cmd.Username, Role: "customer", Status: "active"}, nil
}
func (s *stubRepo) GetCustomerByID(ctx context.Context, customerID uint) (*models.User, error) {
	return &models.User{ID: customerID, Username: "u"}, nil
}
func (s *stubRepo) UpdateCustomer(ctx context.Context, customerID uint, cmd UpdateCustomerCommand) (*models.User, error) {
	s.lastUpdate = cmd
	return &models.User{ID: customerID}, nil
}
func (s *stubRepo) ListCustomers(ctx context.Context, query ListCustomersQuery) ([]CustomerInfoDTO, int64, error) {
	return []CustomerInfoDTO{{User: models.User{ID: 1}}}, 1, nil
}
func (s *stubRepo) GetCustomerActivity(ctx context.Context, customerID uint, limit int) (*CustomerActivityDTO, error) {
	return &CustomerActivityDTO{CustomerID: customerID}, nil
}
func (s *stubRepo) AddNote(ctx context.Context, customerID uint, note CustomerNoteDTO) error {
	s.lastNote = note
	return nil
}
func (s *stubRepo) UpdateTags(ctx context.Context, customerID uint, tags []string) error {
	s.lastTags = tags
	return nil
}
func (s *stubRepo) GetStats(ctx context.Context) (*CustomerStatsDTO, error) {
	return &CustomerStatsDTO{Total: 1}, nil
}

func TestCreateCustomerNormalizesDefaults(t *testing.T) {
	repo := &stubRepo{}
	svc := NewService(repo)
	_, err := svc.CreateCustomer(context.Background(), CreateCustomerCommand{
		Username: "alice",
		Tags:     []string{" vip ", "vip", "beta"},
	})
	if err != nil {
		t.Fatalf("CreateCustomer() error = %v", err)
	}
	if repo.lastCreate.Source != "web" {
		t.Fatalf("expected default source web, got %s", repo.lastCreate.Source)
	}
	if repo.lastCreate.Priority != "normal" {
		t.Fatalf("expected default priority normal, got %s", repo.lastCreate.Priority)
	}
	if len(repo.lastCreate.Tags) != 2 {
		t.Fatalf("expected deduped tags, got %+v", repo.lastCreate.Tags)
	}
}

func TestUpdateTagsNormalizesTags(t *testing.T) {
	repo := &stubRepo{}
	svc := NewService(repo)
	if err := svc.UpdateTags(context.Background(), 1, []string{" a ", "a", "b"}); err != nil {
		t.Fatalf("UpdateTags() error = %v", err)
	}
	if len(repo.lastTags) != 2 {
		t.Fatalf("expected 2 tags, got %+v", repo.lastTags)
	}
}
