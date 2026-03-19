package voiceprotocol

import (
	"fmt"
	"sort"
)

// Registry stores protocol adapters behind normalized interfaces.
type Registry struct {
	signaling map[Protocol]CallSignalingAdapter
	media     map[Protocol]MediaSessionAdapter
}

func NewRegistry() *Registry {
	return &Registry{
		signaling: make(map[Protocol]CallSignalingAdapter),
		media:     make(map[Protocol]MediaSessionAdapter),
	}
}

func (r *Registry) RegisterSignaling(adapter CallSignalingAdapter) error {
	if r == nil {
		return fmt.Errorf("voice protocol registry is nil")
	}
	if adapter == nil {
		return fmt.Errorf("signaling adapter is nil")
	}
	r.signaling[adapter.Protocol()] = adapter
	return nil
}

func (r *Registry) RegisterMedia(adapter MediaSessionAdapter) error {
	if r == nil {
		return fmt.Errorf("voice protocol registry is nil")
	}
	if adapter == nil {
		return fmt.Errorf("media adapter is nil")
	}
	r.media[adapter.Protocol()] = adapter
	return nil
}

func (r *Registry) Signaling(protocol Protocol) (CallSignalingAdapter, bool) {
	if r == nil {
		return nil, false
	}
	adapter, ok := r.signaling[protocol]
	return adapter, ok
}

func (r *Registry) Media(protocol Protocol) (MediaSessionAdapter, bool) {
	if r == nil {
		return nil, false
	}
	adapter, ok := r.media[protocol]
	return adapter, ok
}

func (r *Registry) SupportedProtocols() []Protocol {
	if r == nil {
		return nil
	}
	seen := make(map[Protocol]struct{}, len(r.signaling)+len(r.media))
	for protocol := range r.signaling {
		seen[protocol] = struct{}{}
	}
	for protocol := range r.media {
		seen[protocol] = struct{}{}
	}
	protocols := make([]Protocol, 0, len(seen))
	for protocol := range seen {
		protocols = append(protocols, protocol)
	}
	sort.Slice(protocols, func(i, j int) bool { return protocols[i] < protocols[j] })
	return protocols
}
