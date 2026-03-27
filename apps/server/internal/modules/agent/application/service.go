package application

import (
	"context"
	"fmt"
	"slices"
	"time"

	"servify/apps/server/internal/models"
	agentdomain "servify/apps/server/internal/modules/agent/domain"
)

type Service struct {
	repo     Repository
	registry RuntimeRegistry
}

func NewService(repo Repository, registry RuntimeRegistry) *Service {
	return &Service{repo: repo, registry: registry}
}

func (s *Service) CreateAgent(ctx context.Context, cmd CreateAgentCommand) (*models.Agent, error) {
	if cmd.UserID == 0 {
		return nil, fmt.Errorf("user_id required")
	}
	profile, model, err := s.repo.GetAgentByUserID(ctx, cmd.UserID)
	if err == nil && profile != nil && model != nil {
		return nil, fmt.Errorf("user is already an agent")
	}
	_, err = s.repo.CreateAgent(ctx, cmd.UserID, cmd.Department, sanitizeSkills(cmd.Skills), normalizeChatConcurrency(cmd.MaxChatConcurrency))
	if err != nil {
		return nil, err
	}
	_, model, err = s.repo.GetAgentByUserID(ctx, cmd.UserID)
	return model, err
}

func (s *Service) GetAgentByUserID(ctx context.Context, userID uint) (*models.Agent, error) {
	_, model, err := s.repo.GetAgentByUserID(ctx, userID)
	return model, err
}

func (s *Service) ListAgents(ctx context.Context, limit int) ([]models.Agent, error) {
	return s.repo.ListAgents(ctx, limit)
}

func (s *Service) GoOnline(ctx context.Context, userID uint) error {
	profile, _, err := s.repo.GetAgentByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if err := s.repo.UpdatePresenceStatus(ctx, userID, agentdomain.PresenceStatusOnline); err != nil {
		return err
	}
	_, err = s.registry.GoOnline(*profile)
	return err
}

func (s *Service) GoOffline(ctx context.Context, userID uint) error {
	if err := s.repo.UpdatePresenceStatus(ctx, userID, agentdomain.PresenceStatusOffline); err != nil {
		return err
	}
	s.registry.GoOffline(userID)
	return nil
}

func (s *Service) UpdateStatus(ctx context.Context, userID uint, status string) error {
	next, err := parsePresenceStatus(status)
	if err != nil {
		return err
	}
	if err := s.repo.UpdatePresenceStatus(ctx, userID, next); err != nil {
		return err
	}
	s.registry.UpdateStatus(userID, next)
	return nil
}

func (s *Service) MarkBusy(ctx context.Context, userID uint) error {
	return s.UpdateStatus(ctx, userID, string(agentdomain.PresenceStatusBusy))
}

func (s *Service) MarkAway(ctx context.Context, userID uint) error {
	return s.UpdateStatus(ctx, userID, string(agentdomain.PresenceStatusAway))
}

func (s *Service) AssignSession(ctx context.Context, sessionID string, agentUserID uint) error {
	if _, ok := s.registry.Get(agentUserID); !ok {
		return fmt.Errorf("agent %d is not online", agentUserID)
	}
	session, err := s.repo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return err
	}
	runtime, err := s.registry.AssignSession(agentUserID, session)
	if err != nil {
		return err
	}
	if err := s.repo.AssignSession(ctx, sessionID, agentUserID); err != nil {
		return err
	}
	return s.repo.UpdateChatLoad(ctx, agentUserID, runtime.CurrentChatLoad)
}

func (s *Service) ReleaseSession(ctx context.Context, sessionID string, agentUserID uint) error {
	if err := s.repo.ReleaseSession(ctx, sessionID, agentUserID); err != nil {
		return err
	}
	runtime, ok := s.registry.ReleaseSession(agentUserID, sessionID)
	if ok {
		return s.repo.UpdateChatLoad(ctx, agentUserID, runtime.CurrentChatLoad)
	}
	return nil
}

func (s *Service) FindAvailableAgent(ctx context.Context, skills []string, priority string) (*AgentRuntimeDTO, error) {
	_ = ctx
	requiredSkills := sanitizeSkills(skills)
	var best *AgentRuntimeDTO
	bestScore := -1.0
	for _, candidate := range s.registry.List() {
		if candidate.Status != string(agentdomain.PresenceStatusOnline) {
			continue
		}
		if candidate.CurrentChatLoad >= candidate.MaxChatConcurrency {
			continue
		}
		score := calculateScore(candidate, requiredSkills, priority)
		if score > bestScore {
			copy := candidate
			best = &copy
			bestScore = score
		}
	}
	if best == nil {
		return nil, fmt.Errorf("no available agent found")
	}
	return best, nil
}

func (s *Service) GetOnlineAgents(ctx context.Context) []AgentRuntimeDTO {
	_ = ctx
	return s.registry.List()
}

func (s *Service) GetStats(ctx context.Context, agentUserID *uint) (*AgentStatsDTO, error) {
	stats, err := s.repo.GetStats(ctx, agentUserID)
	if err != nil {
		return nil, err
	}
	for _, candidate := range s.registry.List() {
		stats.Online++
		if candidate.Status == string(agentdomain.PresenceStatusBusy) || candidate.CurrentChatLoad >= candidate.MaxChatConcurrency {
			stats.Busy++
		}
	}
	return stats, nil
}

func (s *Service) RevokeUserTokens(ctx context.Context, userID uint, revokeAt time.Time) (int, error) {
	if userID == 0 {
		return 0, fmt.Errorf("user_id required")
	}
	if revokeAt.IsZero() {
		revokeAt = time.Now().UTC()
	}
	return s.repo.RevokeUserTokens(ctx, userID, revokeAt)
}

func (s *Service) ApplySessionTransfer(ctx context.Context, sessionID string, fromAgentID *uint, toAgentID uint) error {
	_ = ctx
	s.registry.ApplyTransfer(sessionID, fromAgentID, toAgentID)
	if runtime, ok := s.registry.Get(toAgentID); ok {
		if err := s.repo.UpdateChatLoad(ctx, toAgentID, runtime.CurrentChatLoad); err != nil {
			return err
		}
	}
	if fromAgentID != nil {
		if runtime, ok := s.registry.Get(*fromAgentID); ok {
			if err := s.repo.UpdateChatLoad(ctx, *fromAgentID, runtime.CurrentChatLoad); err != nil {
				return err
			}
		}
	}
	return nil
}

func sanitizeSkills(skills []string) []string {
	out := make([]string, 0, len(skills))
	for _, skill := range skills {
		if skill == "" {
			continue
		}
		if !slices.Contains(out, skill) {
			out = append(out, skill)
		}
	}
	return out
}

func normalizeChatConcurrency(v int) int {
	if v <= 0 {
		return 5
	}
	return v
}

func parsePresenceStatus(value string) (agentdomain.PresenceStatus, error) {
	switch agentdomain.PresenceStatus(value) {
	case agentdomain.PresenceStatusOnline, agentdomain.PresenceStatusBusy, agentdomain.PresenceStatusAway, agentdomain.PresenceStatusOffline:
		return agentdomain.PresenceStatus(value), nil
	default:
		return "", fmt.Errorf("invalid status: %s", value)
	}
}

func calculateScore(agent AgentRuntimeDTO, requiredSkills []string, priority string) float64 {
	score := agent.Rating
	if agent.MaxChatConcurrency > 0 {
		loadRatio := float64(agent.CurrentChatLoad) / float64(agent.MaxChatConcurrency)
		score += (1 - loadRatio) * 3
	}
	if agent.AvgResponseTime > 0 {
		responseScore := 300.0 / float64(agent.AvgResponseTime)
		if responseScore > 2 {
			responseScore = 2
		}
		score += responseScore
	}
	if len(requiredSkills) > 0 {
		matched := 0
		for _, required := range requiredSkills {
			if slices.Contains(agent.Skills, required) {
				matched++
			}
		}
		score += float64(matched) / float64(len(requiredSkills)) * 2
	}
	switch priority {
	case "urgent":
		score += 0.2
	case "high":
		score += 0.1
	}
	return score
}
