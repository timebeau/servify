package delivery

import (
	"context"
	"testing"

	"servify/apps/server/internal/platform/voiceprotocol"
)

func TestRTPAdapterMapSessionStarted(t *testing.T) {
	adapter := NewRTPAdapter()
	event, err := adapter.MapSessionStarted(context.Background(), PacketMediaPayload{
		CallID:       "call-rtp-1",
		ConnectionID: "rtp-conn-1",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if event.Protocol != voiceprotocol.ProtocolRTP || event.Kind != voiceprotocol.MediaEventSessionStarted {
		t.Fatalf("unexpected event: %+v", event)
	}
}

func TestSRTPAdapterMapSessionStarted(t *testing.T) {
	adapter := NewSRTPAdapter()
	event, err := adapter.MapSessionStarted(context.Background(), PacketMediaPayload{
		CallID:       "call-srtp-1",
		ConnectionID: "srtp-conn-1",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if event.Protocol != voiceprotocol.ProtocolSRTP || event.Kind != voiceprotocol.MediaEventSessionStarted {
		t.Fatalf("unexpected event: %+v", event)
	}
}
