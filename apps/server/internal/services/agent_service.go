package services

import (
	"context"
	"strings"
	"sync"
	"time"

	"servify/apps/server/internal/models"
	agentapp "servify/apps/server/internal/modules/agent/application"
	agentinfra "servify/apps/server/internal/modules/agent/infra"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// AgentService 人工客服服务兼容层。
type AgentService struct {
	db           *gorm.DB
	logger       *logrus.Logger
	onlineAgents sync.Map // map[uint]*AgentInfo
	agentQueues  sync.Map // map[uint]chan *models.Session
	module       *agentapp.Service
}

// NewAgentService 创建人工客服服务。
func NewAgentService(db *gorm.DB, logger *logrus.Logger) *AgentService {
	if logger == nil {
		logger = logrus.New()
	}

	repo := agentinfra.NewGormRepository(db)
	registry := agentinfra.NewInMemoryRegistry()
	service := &AgentService{
		db:     db,
		logger: logger,
		module: agentapp.NewService(repo, registry),
	}

	go service.backgroundTasks()
	return service
}

// AgentInfo 在线客服信息。
type AgentInfo struct {
	UserID          uint                       `json:"user_id"`
	Username        string                     `json:"username"`
	Name            string                     `json:"name"`
	Department      string                     `json:"department"`
	Skills          []string                   `json:"skills"`
	Status          string                     `json:"status"`
	MaxConcurrent   int                        `json:"max_concurrent"`
	CurrentLoad     int                        `json:"current_load"`
	Rating          float64                    `json:"rating"`
	AvgResponseTime int                        `json:"avg_response_time"`
	LastActivity    time.Time                  `json:"last_activity"`
	ConnectedAt     time.Time                  `json:"connected_at"`
	Sessions        map[string]*models.Session `json:"-"`
}

// AgentCreateRequest 创建客服请求。
type AgentCreateRequest struct {
	UserID        uint   `json:"user_id" binding:"required"`
	Department    string `json:"department"`
	Skills        string `json:"skills"`
	MaxConcurrent int    `json:"max_concurrent"`
}

// AgentUpdateRequest 更新客服请求。
type AgentUpdateRequest struct {
	Department    *string `json:"department"`
	Skills        *string `json:"skills"`
	Status        *string `json:"status"`
	MaxConcurrent *int    `json:"max_concurrent"`
}

func (s *AgentService) CreateAgent(ctx context.Context, req *AgentCreateRequest) (*models.Agent, error) {
	return s.module.CreateAgent(ctx, agentapp.CreateAgentCommand{
		UserID:             req.UserID,
		Department:         req.Department,
		Skills:             s.parseSkills(req.Skills),
		MaxChatConcurrency: req.MaxConcurrent,
	})
}

func (s *AgentService) GetAgentByUserID(ctx context.Context, userID uint) (*models.Agent, error) {
	return s.module.GetAgentByUserID(ctx, userID)
}

func (s *AgentService) ListAgents(ctx context.Context, limit int) ([]models.Agent, error) {
	return s.module.ListAgents(ctx, limit)
}

func (s *AgentService) AgentGoOnline(ctx context.Context, userID uint) error {
	if err := s.module.GoOnline(ctx, userID); err != nil {
		return err
	}
	s.syncLegacyRuntime(ctx)
	s.logger.Infof("Agent %d went online", userID)
	return nil
}

func (s *AgentService) AgentGoOffline(ctx context.Context, userID uint) error {
	if err := s.module.GoOffline(ctx, userID); err != nil {
		return err
	}
	s.syncLegacyRuntime(ctx)
	if queue, ok := s.agentQueues.LoadAndDelete(userID); ok {
		close(queue.(chan *models.Session))
	}
	s.logger.Infof("Agent %d went offline", userID)
	return nil
}

func (s *AgentService) UpdateAgentStatus(ctx context.Context, userID uint, status string) error {
	if err := s.module.UpdateStatus(ctx, userID, status); err != nil {
		return err
	}
	s.syncLegacyRuntime(ctx)
	s.logger.Infof("Agent %d status updated to %s", userID, status)
	return nil
}

func (s *AgentService) AssignSessionToAgent(ctx context.Context, sessionID string, agentID uint) error {
	if err := s.module.AssignSession(ctx, sessionID, agentID); err != nil {
		return err
	}
	s.syncLegacyRuntime(ctx)
	s.logger.Infof("Assigned session %s to agent %d", sessionID, agentID)
	return nil
}

func (s *AgentService) ReleaseSessionFromAgent(ctx context.Context, sessionID string, agentID uint) error {
	if err := s.module.ReleaseSession(ctx, sessionID, agentID); err != nil {
		return err
	}
	s.syncLegacyRuntime(ctx)
	s.logger.Infof("Released session %s from agent %d", sessionID, agentID)
	return nil
}

func (s *AgentService) FindAvailableAgent(ctx context.Context, skills []string, priority string) (*AgentInfo, error) {
	runtime, err := s.module.FindAvailableAgent(ctx, skills, priority)
	if err != nil {
		return nil, err
	}
	return mapRuntimeToLegacy(runtime), nil
}

func (s *AgentService) GetOnlineAgents(ctx context.Context) []*AgentInfo {
	runtimes := s.module.GetOnlineAgents(ctx)
	out := make([]*AgentInfo, 0, len(runtimes))
	for _, item := range runtimes {
		copy := item
		out = append(out, mapRuntimeToLegacy(&copy))
	}
	return out
}

func (s *AgentService) GetAgentStats(ctx context.Context, agentID *uint) (*AgentStats, error) {
	stats, err := s.module.GetStats(ctx, agentID)
	if err != nil {
		return nil, err
	}
	return &AgentStats{
		Total:           stats.Total,
		Online:          stats.Online,
		Busy:            stats.Busy,
		AvgResponseTime: stats.AvgResponseTime,
		AvgRating:       stats.AvgRating,
	}, nil
}

func (s *AgentService) parseSkills(skillsStr string) []string {
	if skillsStr == "" {
		return []string{}
	}
	skills := []string{}
	for _, skill := range strings.Split(skillsStr, ",") {
		skill = strings.TrimSpace(skill)
		if skill != "" {
			skills = append(skills, skill)
		}
	}
	return skills
}

func (s *AgentService) backgroundTasks() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.cleanupInactiveAgents()
		s.updateAgentMetrics()
	}
}

func (s *AgentService) cleanupInactiveAgents() {
	timeout := 5 * time.Minute
	runtimes := s.module.GetOnlineAgents(context.Background())
	for _, item := range runtimes {
		if time.Since(item.LastActivity) > timeout {
			s.logger.Warnf("Agent %d appears inactive, marking as away", item.UserID)
			_ = s.module.MarkAway(context.Background(), item.UserID)
		}
	}
	s.syncLegacyRuntime(context.Background())
}

func (s *AgentService) updateAgentMetrics() {}

func (s *AgentService) ApplySessionTransfer(sessionID string, fromAgentID *uint, toAgentID uint) {
	ctx := context.Background()
	if err := s.module.ApplySessionTransfer(ctx, sessionID, fromAgentID, toAgentID); err != nil {
		s.logger.Warnf("apply session transfer to agent module failed: %v", err)
	}
	s.syncLegacyRuntime(ctx)
}

func (s *AgentService) syncLegacyRuntime(ctx context.Context) {
	runtimes := s.module.GetOnlineAgents(ctx)
	active := make(map[uint]struct{}, len(runtimes))
	for _, runtime := range runtimes {
		active[runtime.UserID] = struct{}{}
		s.onlineAgents.Store(runtime.UserID, mapRuntimeToLegacy(&runtime))
		if _, ok := s.agentQueues.Load(runtime.UserID); !ok {
			queueSize := runtime.MaxChatConcurrency * 2
			if queueSize <= 0 {
				queueSize = 10
			}
			s.agentQueues.Store(runtime.UserID, make(chan *models.Session, queueSize))
		}
	}
	var stale []uint
	s.onlineAgents.Range(func(key, value any) bool {
		userID, ok := key.(uint)
		if ok {
			if _, exists := active[userID]; !exists {
				stale = append(stale, userID)
			}
		}
		return true
	})
	for _, userID := range stale {
		s.onlineAgents.Delete(userID)
		if queue, ok := s.agentQueues.LoadAndDelete(userID); ok {
			close(queue.(chan *models.Session))
		}
	}
}

func mapRuntimeToLegacy(runtime *agentapp.AgentRuntimeDTO) *AgentInfo {
	if runtime == nil {
		return nil
	}
	return &AgentInfo{
		UserID:          runtime.UserID,
		Username:        runtime.Username,
		Name:            runtime.Name,
		Department:      runtime.Department,
		Skills:          append([]string(nil), runtime.Skills...),
		Status:          runtime.Status,
		MaxConcurrent:   runtime.MaxChatConcurrency,
		CurrentLoad:     runtime.CurrentChatLoad,
		Rating:          runtime.Rating,
		AvgResponseTime: runtime.AvgResponseTime,
		LastActivity:    runtime.LastActivity,
		ConnectedAt:     runtime.ConnectedAt,
		Sessions:        make(map[string]*models.Session),
	}
}

// AgentStats 客服统计信息。
type AgentStats struct {
	Total           int64   `json:"total"`
	Online          int64   `json:"online"`
	Busy            int64   `json:"busy"`
	AvgResponseTime int64   `json:"avg_response_time"`
	AvgRating       float64 `json:"avg_rating"`
}
