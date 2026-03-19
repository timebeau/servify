package pstnprovider

import (
	"context"
	"testing"

	"servify/apps/server/internal/platform/voiceprotocol"
)

func TestAdapterMapHangup(t *testing.T) {
	adapter := NewAdapter()
	event, err := adapter.MapHangup(context.Background(), WebhookEvent{
		Provider:  "twilio",
		EventType: "call.completed",
		CallID:    "call-9",
		From:      "+1001",
		To:        "+2002",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if event.Protocol != voiceprotocol.ProtocolPSTNProvider || event.Kind != voiceprotocol.CallEventHangup {
		t.Fatalf("unexpected event: %+v", event)
	}
	if event.Metadata["provider"] != "twilio" {
		t.Fatalf("expected provider metadata, got %+v", event.Metadata)
	}
}
