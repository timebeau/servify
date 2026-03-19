package sipws

import (
	"context"
	"testing"

	"servify/apps/server/internal/platform/voiceprotocol"
)

func TestAdapterMapInvite(t *testing.T) {
	adapter := NewAdapter()
	event, err := adapter.MapInvite(context.Background(), SignalingMessage{
		CallID: "call-1",
		Method: "INVITE",
		From:   "alice",
		To:     "bob",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if event.Protocol != voiceprotocol.ProtocolSIPWebSocket || event.Kind != voiceprotocol.CallEventInvite {
		t.Fatalf("unexpected event: %+v", event)
	}
}
