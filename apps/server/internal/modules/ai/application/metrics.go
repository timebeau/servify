package application

import (
	"sync"
	"time"
)

// Metrics captures high-level orchestrator metrics.
type Metrics struct {
	mu                 sync.RWMutex
	queryCount         int64
	successCount       int64
	fallbackCount      int64
	providerUsageCount map[string]int64
	lastLatency        time.Duration
	lastTokenUsage     int
}

func NewMetrics() *Metrics {
	return &Metrics{
		providerUsageCount: make(map[string]int64),
	}
}

func (m *Metrics) RecordQuery() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.queryCount++
}

func (m *Metrics) RecordSuccess(provider string, latency time.Duration, totalTokens int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.successCount++
	m.lastLatency = latency
	m.lastTokenUsage = totalTokens
	if provider != "" {
		m.providerUsageCount[provider]++
	}
}

func (m *Metrics) RecordFallback() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.fallbackCount++
}

type MetricsSnapshot struct {
	QueryCount         int64            `json:"query_count"`
	SuccessCount       int64            `json:"success_count"`
	FallbackCount      int64            `json:"fallback_count"`
	ProviderUsageCount map[string]int64 `json:"provider_usage_count"`
	LastLatency        time.Duration    `json:"last_latency"`
	LastTokenUsage     int              `json:"last_token_usage"`
}

func (m *Metrics) Snapshot() MetricsSnapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := MetricsSnapshot{
		QueryCount:         m.queryCount,
		SuccessCount:       m.successCount,
		FallbackCount:      m.fallbackCount,
		ProviderUsageCount: make(map[string]int64, len(m.providerUsageCount)),
		LastLatency:        m.lastLatency,
		LastTokenUsage:     m.lastTokenUsage,
	}
	for k, v := range m.providerUsageCount {
		out.ProviderUsageCount[k] = v
	}
	return out
}
