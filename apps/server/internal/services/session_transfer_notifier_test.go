package services

import (
	"testing"
	"time"
)

type stubSessionTransferRealtimeSink struct {
	sessionID string
	message   WebSocketMessage
}

func (s *stubSessionTransferRealtimeSink) SendToSession(sessionID string, message WebSocketMessage) {
	s.sessionID = sessionID
	s.message = message
}

func TestWebsocketSessionTransferNotifierNotifyTransfer(t *testing.T) {
	sink := &stubSessionTransferRealtimeSink{}
	notifier := NewSessionTransferNotifier(sink)
	at := time.Unix(1700000000, 0)

	notifier.NotifyTransfer("s1", 9, "transferred", at)

	if sink.sessionID != "s1" {
		t.Fatalf("sessionID = %q", sink.sessionID)
	}
	if sink.message.Type != "transfer_notification" {
		t.Fatalf("type = %q", sink.message.Type)
	}
	data, ok := sink.message.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("data type = %T", sink.message.Data)
	}
	if data["message"] != "transferred" || data["agent_id"] != uint(9) || data["timestamp"] != at {
		t.Fatalf("unexpected data: %+v", data)
	}
}

func TestWebsocketSessionTransferNotifierNotifyWaiting(t *testing.T) {
	sink := &stubSessionTransferRealtimeSink{}
	notifier := NewSessionTransferNotifier(sink)
	at := time.Unix(1700000001, 0)

	notifier.NotifyWaiting("s2", "queued", at)

	if sink.sessionID != "s2" {
		t.Fatalf("sessionID = %q", sink.sessionID)
	}
	if sink.message.Type != "waiting_notification" {
		t.Fatalf("type = %q", sink.message.Type)
	}
	data, ok := sink.message.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("data type = %T", sink.message.Data)
	}
	if data["message"] != "queued" || data["timestamp"] != at {
		t.Fatalf("unexpected data: %+v", data)
	}
}

func TestNewSessionTransferNotifierNilSink(t *testing.T) {
	if got := NewSessionTransferNotifier(nil); got != nil {
		t.Fatalf("expected nil notifier, got %#v", got)
	}
}

