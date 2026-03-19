package channel

import "testing"

func TestNewWebMessageEvent(t *testing.T) {
	event := NewWebMessageEvent("conv-1", "customer-1", "hello")

	if event.Channel != WebChannel {
		t.Fatalf("expected web channel, got %+v", event)
	}
	if event.Kind != EventKindMessage {
		t.Fatalf("expected message kind, got %+v", event)
	}
	if event.Payload["content"] != "hello" {
		t.Fatalf("expected content payload, got %+v", event.Payload)
	}
	if event.EventID == "" {
		t.Fatalf("expected generated event id")
	}
}
