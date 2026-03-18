package application

import (
	"context"
	"testing"

	"servify/apps/server/internal/models"
	agentdomain "servify/apps/server/internal/modules/agent/domain"

	"gorm.io/gorm"
)

type stubRepo struct {
	profile         *agentdomain.AgentProfile
	model           *models.Agent
	stats           *AgentStatsDTO
	statusUpdates   []string
	lastChatLoad    int
	sessionAssigned string
}

func (s *stubRepo) CreateAgent(ctx context.Context, userID uint, department string, skills []string, maxChatConcurrency int) (*agentdomain.AgentProfile, error) {
	s.profile = &agentdomain.AgentProfile{
		UserID:              userID,
		Department:          department,
		Skills:              append([]string(nil), skills...),
		MaxChatConcurrency:  maxChatConcurrency,
		MaxVoiceConcurrency: 1,
		Rating:              5,
	}
	s.model = &models.Agent{UserID: userID, Department: department, MaxConcurrent: maxChatConcurrency}
	return s.profile, nil
}

func (s *stubRepo) GetAgentByUserID(ctx context.Context, userID uint) (*agentdomain.AgentProfile, *models.Agent, error) {
	if s.profile == nil || s.model == nil || s.model.UserID != userID {
		return nil, nil, gorm.ErrRecordNotFound
	}
	return s.profile, s.model, nil
}

func (s *stubRepo) ListAgents(ctx context.Context, limit int) ([]models.Agent, error) {
	if s.model == nil {
		return nil, nil
	}
	return []models.Agent{*s.model}, nil
}

func (s *stubRepo) UpdatePresenceStatus(ctx context.Context, userID uint, status agentdomain.PresenceStatus) error {
	s.statusUpdates = append(s.statusUpdates, string(status))
	return nil
}

func (s *stubRepo) UpdateChatLoad(ctx context.Context, userID uint, currentLoad int) error {
	s.lastChatLoad = currentLoad
	return nil
}

func (s *stubRepo) GetSessionByID(ctx context.Context, sessionID string) (*models.Session, error) {
	return &models.Session{ID: sessionID}, nil
}

func (s *stubRepo) AssignSession(ctx context.Context, sessionID string, agentUserID uint) error {
	s.sessionAssigned = sessionID
	return nil
}

func (s *stubRepo) ReleaseSession(ctx context.Context, sessionID string, agentUserID uint) error {
	return nil
}

func (s *stubRepo) GetStats(ctx context.Context, agentUserID *uint) (*AgentStatsDTO, error) {
	if s.stats == nil {
		s.stats = &AgentStatsDTO{Total: 1}
	}
	return s.stats, nil
}

func TestServiceGoOnlineAndAssignSession(t *testing.T) {
	repo := &stubRepo{
		profile: &agentdomain.AgentProfile{
			UserID:              7,
			Username:            "agent7",
			Name:                "Agent Seven",
			Skills:              []string{"billing", "chat"},
			MaxChatConcurrency:  2,
			MaxVoiceConcurrency: 1,
			Rating:              4.8,
		},
		model: &models.Agent{UserID: 7, MaxConcurrent: 2},
	}
	registry := newStubRegistry()
	svc := NewService(repo, registry)

	if err := svc.GoOnline(context.Background(), 7); err != nil {
		t.Fatalf("GoOnline() error = %v", err)
	}
	if err := svc.AssignSession(context.Background(), "sess-1", 7); err != nil {
		t.Fatalf("AssignSession() error = %v", err)
	}
	got, err := svc.FindAvailableAgent(context.Background(), []string{"billing"}, "high")
	if err != nil {
		t.Fatalf("FindAvailableAgent() error = %v", err)
	}
	if got.UserID != 7 {
		t.Fatalf("FindAvailableAgent() user_id = %d, want 7", got.UserID)
	}
	if repo.lastChatLoad != 1 {
		t.Fatalf("UpdateChatLoad() = %d, want 1", repo.lastChatLoad)
	}
}

type stubRegistry struct {
	items map[uint]AgentRuntimeDTO
}

func newStubRegistry() *stubRegistry {
	return &stubRegistry{items: make(map[uint]AgentRuntimeDTO)}
}

func (s *stubRegistry) GoOnline(profile agentdomain.AgentProfile) (AgentRuntimeDTO, error) {
	item := AgentRuntimeDTO{
		UserID:              profile.UserID,
		Username:            profile.Username,
		Name:                profile.Name,
		Department:          profile.Department,
		Skills:              append([]string(nil), profile.Skills...),
		Status:              "online",
		MaxChatConcurrency:  profile.MaxChatConcurrency,
		MaxVoiceConcurrency: profile.MaxVoiceConcurrency,
		CurrentChatLoad:     profile.CurrentChatLoad,
		CurrentVoiceLoad:    profile.CurrentVoiceLoad,
		Rating:              profile.Rating,
		AvgResponseTime:     profile.AvgResponseTime,
	}
	s.items[profile.UserID] = item
	return item, nil
}

func (s *stubRegistry) GoOffline(userID uint) {
	delete(s.items, userID)
}

func (s *stubRegistry) UpdateStatus(userID uint, status agentdomain.PresenceStatus) {
	item := s.items[userID]
	item.Status = string(status)
	s.items[userID] = item
}

func (s *stubRegistry) AssignSession(userID uint, session *models.Session) (AgentRuntimeDTO, error) {
	item := s.items[userID]
	item.CurrentChatLoad++
	s.items[userID] = item
	return item, nil
}

func (s *stubRegistry) ReleaseSession(userID uint, sessionID string) (AgentRuntimeDTO, bool) {
	item, ok := s.items[userID]
	if !ok {
		return AgentRuntimeDTO{}, false
	}
	if item.CurrentChatLoad > 0 {
		item.CurrentChatLoad--
	}
	s.items[userID] = item
	return item, true
}

func (s *stubRegistry) ApplyTransfer(sessionID string, fromAgentID *uint, toAgentID uint) {}

func (s *stubRegistry) Get(userID uint) (AgentRuntimeDTO, bool) {
	item, ok := s.items[userID]
	return item, ok
}

func (s *stubRegistry) List() []AgentRuntimeDTO {
	out := make([]AgentRuntimeDTO, 0, len(s.items))
	for _, item := range s.items {
		out = append(out, item)
	}
	return out
}
