package infra

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"servify/apps/server/internal/models"
	agentapp "servify/apps/server/internal/modules/agent/application"
	agentdomain "servify/apps/server/internal/modules/agent/domain"
)

type runtimeAgent struct {
	agentapp.AgentRuntimeDTO
	Sessions map[string]*models.Session
}

type InMemoryRegistry struct {
	mu     sync.RWMutex
	agents map[uint]*runtimeAgent
}

func NewInMemoryRegistry() *InMemoryRegistry {
	return &InMemoryRegistry{agents: make(map[uint]*runtimeAgent)}
}

func (r *InMemoryRegistry) GoOnline(profile agentdomain.AgentProfile) (agentapp.AgentRuntimeDTO, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	item := &runtimeAgent{
		AgentRuntimeDTO: agentapp.AgentRuntimeDTO{
			UserID:              profile.UserID,
			Username:            profile.Username,
			Name:                profile.Name,
			Department:          profile.Department,
			Skills:              append([]string(nil), profile.Skills...),
			Status:              string(agentdomain.PresenceStatusOnline),
			MaxChatConcurrency:  profile.MaxChatConcurrency,
			MaxVoiceConcurrency: profile.MaxVoiceConcurrency,
			CurrentChatLoad:     profile.CurrentChatLoad,
			CurrentVoiceLoad:    profile.CurrentVoiceLoad,
			Rating:              profile.Rating,
			AvgResponseTime:     profile.AvgResponseTime,
			LastActivity:        now,
			ConnectedAt:         now,
		},
		Sessions: make(map[string]*models.Session),
	}
	r.agents[profile.UserID] = item
	return item.AgentRuntimeDTO, nil
}

func (r *InMemoryRegistry) GoOffline(userID uint) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.agents, userID)
}

func (r *InMemoryRegistry) UpdateStatus(userID uint, status agentdomain.PresenceStatus) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if item, ok := r.agents[userID]; ok {
		item.Status = string(status)
		item.LastActivity = time.Now()
	}
}

func (r *InMemoryRegistry) AssignSession(userID uint, session *models.Session) (agentapp.AgentRuntimeDTO, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	item, ok := r.agents[userID]
	if !ok {
		return agentapp.AgentRuntimeDTO{}, fmt.Errorf("agent %d is not online", userID)
	}
	if item.CurrentChatLoad >= item.MaxChatConcurrency {
		return agentapp.AgentRuntimeDTO{}, fmt.Errorf("agent %d is at maximum capacity", userID)
	}
	if item.Status == string(agentdomain.PresenceStatusOffline) {
		return agentapp.AgentRuntimeDTO{}, fmt.Errorf("agent %d is offline", userID)
	}
	item.CurrentChatLoad++
	item.LastActivity = time.Now()
	copySession := *session
	item.Sessions[session.ID] = &copySession
	return item.AgentRuntimeDTO, nil
}

func (r *InMemoryRegistry) ReleaseSession(userID uint, sessionID string) (agentapp.AgentRuntimeDTO, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	item, ok := r.agents[userID]
	if !ok {
		return agentapp.AgentRuntimeDTO{}, false
	}
	if item.CurrentChatLoad > 0 {
		item.CurrentChatLoad--
	}
	delete(item.Sessions, sessionID)
	item.LastActivity = time.Now()
	return item.AgentRuntimeDTO, true
}

func (r *InMemoryRegistry) ApplyTransfer(sessionID string, fromAgentID *uint, toAgentID uint) {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	if fromAgentID != nil {
		if item, ok := r.agents[*fromAgentID]; ok {
			if item.CurrentChatLoad > 0 {
				item.CurrentChatLoad--
			}
			delete(item.Sessions, sessionID)
			item.LastActivity = now
		}
	}
	if item, ok := r.agents[toAgentID]; ok {
		item.CurrentChatLoad++
		item.Sessions[sessionID] = &models.Session{ID: sessionID, AgentID: &toAgentID, Status: "active"}
		item.LastActivity = now
	}
}

func (r *InMemoryRegistry) Get(userID uint) (agentapp.AgentRuntimeDTO, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, ok := r.agents[userID]
	if !ok {
		return agentapp.AgentRuntimeDTO{}, false
	}
	return item.AgentRuntimeDTO, true
}

func (r *InMemoryRegistry) List() []agentapp.AgentRuntimeDTO {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]agentapp.AgentRuntimeDTO, 0, len(r.agents))
	for _, item := range r.agents {
		out = append(out, item.AgentRuntimeDTO)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].UserID < out[j].UserID })
	return out
}
