package sip

import (
	"context"
	"testing"

	"servify/apps/server/internal/platform/voiceprotocol"
)

func TestVoiceProtocolAdapterMapInvite(t *testing.T) {
	adapter := NewVoiceProtocolAdapter()

	event, err := adapter.MapInvite(context.Background(), InboundCall{
		CallID:         "call-1",
		ConversationID: "conv-1",
		From:           "1001",
		To:             "2001",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if event.Protocol != voiceprotocol.ProtocolSIP || event.Kind != voiceprotocol.CallEventInvite {
		t.Fatalf("unexpected event: %+v", event)
	}
}

func TestVoiceProtocolAdapterRejectsUnsupportedPayload(t *testing.T) {
	adapter := NewVoiceProtocolAdapter()

	if _, err := adapter.MapHangup(context.Background(), "invalid"); err == nil {
		t.Fatalf("expected payload type error")
	}
}
