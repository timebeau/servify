package sip

import (
	"context"
	"testing"
	"time"

	channelplatform "servify/apps/server/internal/platform/channel"
)

func TestDefaultAdapterMapInvite(t *testing.T) {
	adapter := NewDefaultAdapter()
	now := time.Now().UTC()

	event, err := adapter.MapInvite(context.Background(), InboundCall{
		CallID:         "call-1",
		ConversationID: "conv-1",
		From:           "1001",
		To:             "2001",
		Event:          CallEventInvite,
		OccurredAt:     now,
		Metadata:       map[string]interface{}{"tenant_id": "t-1"},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if event.Channel != ChannelName || event.Kind != channelplatform.EventKindCall {
		t.Fatalf("unexpected event: %+v", event)
	}
	if event.Payload["event"] != CallEventInvite {
		t.Fatalf("expected invite payload, got %+v", event.Payload)
	}
}

func TestDefaultAdapterMapDTMF(t *testing.T) {
	adapter := NewDefaultAdapter()

	event, err := adapter.MapDTMF(context.Background(), InboundCall{
		CallID: "call-2",
		From:   "1002",
		To:     "ivr",
		DTMF:   "9",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if event.Payload["digits"] != "9" {
		t.Fatalf("expected dtmf digits, got %+v", event.Payload)
	}
	if event.ActorID != "1002" {
		t.Fatalf("expected actor id to be caller, got %+v", event)
	}
}
