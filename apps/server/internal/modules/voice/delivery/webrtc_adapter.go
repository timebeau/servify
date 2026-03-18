package delivery

import (
	"context"

	voiceapp "servify/apps/server/internal/modules/voice/application"
)

type WebRTCAdapter struct {
	service *voiceapp.Service
}

func NewWebRTCAdapter(service *voiceapp.Service) *WebRTCAdapter {
	return &WebRTCAdapter{service: service}
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
