package delivery

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"servify/apps/server/internal/platform/voiceprotocol"
)

type RTPAdapter struct{}

func NewRTPAdapter() *RTPAdapter {
	return &RTPAdapter{}
}

func (a *RTPAdapter) Name() string {
	return "rtp"
}

func (a *RTPAdapter) Protocol() voiceprotocol.Protocol {
	return voiceprotocol.ProtocolRTP
}

func (a *RTPAdapter) MapSessionStarted(_ context.Context, payload interface{}) (voiceprotocol.MediaEvent, error) {
	return mapPacketMediaEvent(payload, voiceprotocol.ProtocolRTP, voiceprotocol.MediaEventSessionStarted)
}

func (a *RTPAdapter) MapSessionClosed(_ context.Context, payload interface{}) (voiceprotocol.MediaEvent, error) {
	return mapPacketMediaEvent(payload, voiceprotocol.ProtocolRTP, voiceprotocol.MediaEventSessionClosed)
}

func (a *RTPAdapter) MapTrackMuted(_ context.Context, payload interface{}) (voiceprotocol.MediaEvent, error) {
	return mapPacketMediaEvent(payload, voiceprotocol.ProtocolRTP, voiceprotocol.MediaEventTrackMuted)
}

func (a *RTPAdapter) MapTrackUnmuted(_ context.Context, payload interface{}) (voiceprotocol.MediaEvent, error) {
	return mapPacketMediaEvent(payload, voiceprotocol.ProtocolRTP, voiceprotocol.MediaEventTrackUnmuted)
}

func (a *RTPAdapter) MapRecordingStarted(_ context.Context, payload interface{}) (voiceprotocol.MediaEvent, error) {
	return mapPacketMediaEvent(payload, voiceprotocol.ProtocolRTP, voiceprotocol.MediaEventRecordingStart)
}

func (a *RTPAdapter) MapRecordingStopped(_ context.Context, payload interface{}) (voiceprotocol.MediaEvent, error) {
	return mapPacketMediaEvent(payload, voiceprotocol.ProtocolRTP, voiceprotocol.MediaEventRecordingStop)
}

type SRTPAdapter struct{}

func NewSRTPAdapter() *SRTPAdapter {
	return &SRTPAdapter{}
}

func (a *SRTPAdapter) Name() string {
	return "srtp"
}

func (a *SRTPAdapter) Protocol() voiceprotocol.Protocol {
	return voiceprotocol.ProtocolSRTP
}

func (a *SRTPAdapter) MapSessionStarted(_ context.Context, payload interface{}) (voiceprotocol.MediaEvent, error) {
	return mapPacketMediaEvent(payload, voiceprotocol.ProtocolSRTP, voiceprotocol.MediaEventSessionStarted)
}

func (a *SRTPAdapter) MapSessionClosed(_ context.Context, payload interface{}) (voiceprotocol.MediaEvent, error) {
	return mapPacketMediaEvent(payload, voiceprotocol.ProtocolSRTP, voiceprotocol.MediaEventSessionClosed)
}

func (a *SRTPAdapter) MapTrackMuted(_ context.Context, payload interface{}) (voiceprotocol.MediaEvent, error) {
	return mapPacketMediaEvent(payload, voiceprotocol.ProtocolSRTP, voiceprotocol.MediaEventTrackMuted)
}

func (a *SRTPAdapter) MapTrackUnmuted(_ context.Context, payload interface{}) (voiceprotocol.MediaEvent, error) {
	return mapPacketMediaEvent(payload, voiceprotocol.ProtocolSRTP, voiceprotocol.MediaEventTrackUnmuted)
}

func (a *SRTPAdapter) MapRecordingStarted(_ context.Context, payload interface{}) (voiceprotocol.MediaEvent, error) {
	return mapPacketMediaEvent(payload, voiceprotocol.ProtocolSRTP, voiceprotocol.MediaEventRecordingStart)
}

func (a *SRTPAdapter) MapRecordingStopped(_ context.Context, payload interface{}) (voiceprotocol.MediaEvent, error) {
	return mapPacketMediaEvent(payload, voiceprotocol.ProtocolSRTP, voiceprotocol.MediaEventRecordingStop)
}

type PacketMediaPayload struct {
	CallID         string                 `json:"call_id"`
	ConversationID string                 `json:"conversation_id,omitempty"`
	ConnectionID   string                 `json:"connection_id,omitempty"`
	OccurredAt     time.Time              `json:"occurred_at"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

func mapPacketMediaEvent(payload interface{}, protocol voiceprotocol.Protocol, kind voiceprotocol.MediaEventKind) (voiceprotocol.MediaEvent, error) {
	data, err := asPacketMediaPayload(payload)
	if err != nil {
		return voiceprotocol.MediaEvent{}, err
	}
	occurredAt := data.OccurredAt
	if occurredAt.IsZero() {
		occurredAt = time.Now()
	}
	return voiceprotocol.MediaEvent{
		EventID:        fmt.Sprintf("%s-%s-%s", protocol, kind, data.ConnectionID),
		Protocol:       protocol,
		Kind:           kind,
		CallID:         data.CallID,
		ConversationID: data.ConversationID,
		ConnectionID:   data.ConnectionID,
		OccurredAt:     occurredAt,
		Metadata:       data.Metadata,
	}, nil
}

func asPacketMediaPayload(payload interface{}) (PacketMediaPayload, error) {
	switch v := payload.(type) {
	case PacketMediaPayload:
		return v, nil
	case map[string]interface{}:
		raw, err := json.Marshal(v)
		if err != nil {
			return PacketMediaPayload{}, fmt.Errorf("marshal packet media payload: %w", err)
		}
		var data PacketMediaPayload
		if err := json.Unmarshal(raw, &data); err != nil {
			return PacketMediaPayload{}, fmt.Errorf("decode packet media payload: %w", err)
		}
		return data, nil
	default:
		return PacketMediaPayload{}, fmt.Errorf("unsupported packet media payload type %T", payload)
	}
}

var _ voiceprotocol.MediaSessionAdapter = (*RTPAdapter)(nil)
var _ voiceprotocol.MediaSessionAdapter = (*SRTPAdapter)(nil)
