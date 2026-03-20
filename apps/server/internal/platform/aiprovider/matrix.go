package aiprovider

import (
	"slices"
	"time"
)

type Kind string

const (
	KindLLM       Kind = "llm"
	KindKnowledge Kind = "knowledge"
)

type CapabilityName string

const (
	CapabilityChat        CapabilityName = "chat"
	CapabilityChatStream  CapabilityName = "chat_stream"
	CapabilityEmbeddings  CapabilityName = "embeddings"
	CapabilityToolCalling CapabilityName = "tool_calling"
	CapabilityRetrieval   CapabilityName = "retrieval"
	CapabilityIndexing    CapabilityName = "indexing"
	CapabilityDeletion    CapabilityName = "deletion"
	CapabilityHealthCheck CapabilityName = "health_check"
)

type CapabilityDeclaration struct {
	Name     CapabilityName     `json:"name"`
	Enabled  bool               `json:"enabled"`
	Metadata map[string]any     `json:"metadata,omitempty"`
}

type CircuitBreakerPolicy struct {
	Enabled      bool          `json:"enabled"`
	MaxFailures  int           `json:"max_failures,omitempty"`
	ResetTimeout time.Duration `json:"reset_timeout,omitempty"`
}

type FallbackPolicy struct {
	Priority       int                  `json:"priority"`
	FallbackTo     []string             `json:"fallback_to,omitempty"`
	CircuitBreaker CircuitBreakerPolicy `json:"circuit_breaker,omitempty"`
}

type ProviderDescriptor struct {
	ID           string                  `json:"id"`
	Kind         Kind                    `json:"kind"`
	Driver       string                  `json:"driver"`
	Enabled      bool                    `json:"enabled"`
	Capabilities []CapabilityDeclaration `json:"capabilities,omitempty"`
	Fallback     FallbackPolicy          `json:"fallback"`
}

type Matrix struct {
	Providers []ProviderDescriptor `json:"providers"`
}

func (m Matrix) Enabled(kind Kind) []ProviderDescriptor {
	out := make([]ProviderDescriptor, 0, len(m.Providers))
	for _, provider := range m.Providers {
		if provider.Kind == kind && provider.Enabled {
			out = append(out, provider)
		}
	}

	slices.SortStableFunc(out, func(a, b ProviderDescriptor) int {
		if a.Fallback.Priority == b.Fallback.Priority {
			switch {
			case a.ID < b.ID:
				return -1
			case a.ID > b.ID:
				return 1
			default:
				return 0
			}
		}
		if a.Fallback.Priority < b.Fallback.Priority {
			return -1
		}
		return 1
	})

	return out
}

func (m Matrix) Find(id string) (ProviderDescriptor, bool) {
	for _, provider := range m.Providers {
		if provider.ID == id {
			return provider, true
		}
	}
	return ProviderDescriptor{}, false
}
