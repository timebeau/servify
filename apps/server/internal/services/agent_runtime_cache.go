package services

import "sync"

type agentRuntimeCache struct {
	agents sync.Map // map[uint]*AgentInfo
}

func (c *agentRuntimeCache) Store(userID uint, info *AgentInfo) {
	c.agents.Store(userID, info)
}

func (c *agentRuntimeCache) Load(userID uint) (*AgentInfo, bool) {
	value, ok := c.agents.Load(userID)
	if !ok {
		return nil, false
	}
	info, ok := value.(*AgentInfo)
	if !ok {
		return nil, false
	}
	return info, true
}

func (c *agentRuntimeCache) Delete(userID uint) {
	c.agents.Delete(userID)
}

func (c *agentRuntimeCache) CollectStale(active map[uint]struct{}) []uint {
	var stale []uint
	c.agents.Range(func(key, value any) bool {
		userID, ok := key.(uint)
		if ok {
			if _, exists := active[userID]; !exists {
				stale = append(stale, userID)
			}
		}
		return true
	})
	return stale
}
