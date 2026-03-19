package delivery

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	voiceapp "servify/apps/server/internal/modules/voice/application"
	"servify/apps/server/internal/platform/voiceprotocol"
)

type WebRTCAdapter struct {
	service *voiceapp.Service
}

func NewWebRTCAdapter(service *voiceapp.Service) *WebRTCAdapter {
	return &WebRTCAdapter{service: service}
}

func (a *WebRTCAdapter) Name() string {
	return "webrtc"
}

func (a *WebRTCAdapter) Protocol() voiceprotocol.Protocol {
	return voiceprotocol.ProtocolWebRTC
}

func (a *WebRTCAdapter) StartCall(ctx context.Context, sessionID string, connectionID string) {
	if a == nil || a.service == nil {
		return
	}
	_, _ = a.service.StartCall(ctx, voiceapp.StartCallCommand{
		CallID:       connectionID,
		SessionID:    sessionID,
		ConnectionID: connectionID,
	})
}

func (a *WebRTCAdapter) AnswerCall(ctx context.Context, connectionID string) {
	if a == nil || a.service == nil {
		return
	}
	_, _ = a.service.AnswerCall(ctx, voiceapp.AnswerCallCommand{CallID: connectionID})
}

func (a *WebRTCAdapter) EndCall(ctx context.Context, connectionID string) {
	if a == nil || a.service == nil {
		return
	}
	_, _ = a.service.EndCall(ctx, voiceapp.EndCallCommand{CallID: connectionID})
}

func (a *WebRTCAdapter) MapSessionStarted(_ context.Context, payload interface{}) (voiceprotocol.MediaEvent, error) {
	return a.mapMediaEvent(payload, voiceprotocol.MediaEventSessionStarted)
}

func (a *WebRTCAdapter) MapSessionClosed(_ context.Context, payload interface{}) (voiceprotocol.MediaEvent, error) {
	return a.mapMediaEvent(payload, voiceprotocol.MediaEventSessionClosed)
}

func (a *WebRTCAdapter) MapTrackMuted(_ context.Context, payload interface{}) (voiceprotocol.MediaEvent, error) {
	return a.mapMediaEvent(payload, voiceprotocol.MediaEventTrackMuted)
}

func (a *WebRTCAdapter) MapTrackUnmuted(_ context.Context, payload interface{}) (voiceprotocol.MediaEvent, error) {
	return a.mapMediaEvent(payload, voiceprotocol.MediaEventTrackUnmuted)
}

func (a *WebRTCAdapter) MapRecordingStarted(_ context.Context, payload interface{}) (voiceprotocol.MediaEvent, error) {
	return a.mapMediaEvent(payload, voiceprotocol.MediaEventRecordingStart)
}

func (a *WebRTCAdapter) MapRecordingStopped(_ context.Context, payload interface{}) (voiceprotocol.MediaEvent, error) {
	return a.mapMediaEvent(payload, voiceprotocol.MediaEventRecordingStop)
}

func (a *WebRTCAdapter) mapMediaEvent(payload interface{}, kind voiceprotocol.MediaEventKind) (voiceprotocol.MediaEvent, error) {
	data, err := asWebRTCMediaPayload(payload)
	if err != nil {
		return voiceprotocol.MediaEvent{}, err
	}
	occurredAt := data.OccurredAt
	if occurredAt.IsZero() {
		occurredAt = time.Now()
	}
	return voiceprotocol.MediaEvent{
		EventID:        fmt.Sprintf("webrtc-%s-%s", kind, data.ConnectionID),
		Protocol:       voiceprotocol.ProtocolWebRTC,
		Kind:           kind,
		CallID:         data.CallID,
		ConversationID: data.ConversationID,
		ConnectionID:   data.ConnectionID,
		OccurredAt:     occurredAt,
		Metadata:       data.Metadata,
	}, nil
}

func asWebRTCMediaPayload(payload interface{}) (WebRTCMediaPayload, error) {
	switch v := payload.(type) {
	case WebRTCMediaPayload:
		return v, nil
	case map[string]interface{}:
		raw, err := json.Marshal(v)
		if err != nil {
			return WebRTCMediaPayload{}, fmt.Errorf("marshal WebRTC media payload: %w", err)
		}
		var data WebRTCMediaPayload
		if err := json.Unmarshal(raw, &data); err != nil {
			return WebRTCMediaPayload{}, fmt.Errorf("decode WebRTC media payload: %w", err)
		}
		return data, nil
	default:
		return WebRTCMediaPayload{}, fmt.Errorf("unsupported WebRTC payload type %T", payload)
	}
}

var _ voiceprotocol.MediaSessionAdapter = (*WebRTCAdapter)(nil)
