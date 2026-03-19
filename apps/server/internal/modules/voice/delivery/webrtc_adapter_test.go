package delivery

import (
	"context"
	"testing"

	"servify/apps/server/internal/platform/voiceprotocol"
)

func TestWebRTCAdapterMapSessionStarted(t *testing.T) {
	adapter := NewWebRTCAdapter(nil)

	event, err := adapter.MapSessionStarted(context.Background(), WebRTCMediaPayload{
		CallID:         "call-1",
		ConversationID: "conv-1",
		ConnectionID:   "peer-1",
		Metadata:       map[string]interface{}{"track": "audio"},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if event.Protocol != voiceprotocol.ProtocolWebRTC || event.Kind != voiceprotocol.MediaEventSessionStarted {
		t.Fatalf("unexpected event: %+v", event)
	}
}

func TestWebRTCAdapterRejectsUnsupportedPayload(t *testing.T) {
	adapter := NewWebRTCAdapter(nil)

	if _, err := adapter.MapTrackMuted(context.Background(), "invalid"); err == nil {
		t.Fatalf("expected payload type error")
	}
}
