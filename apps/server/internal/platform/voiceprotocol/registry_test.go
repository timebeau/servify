package voiceprotocol

import (
	"context"
	"testing"
)

type testSignalingAdapter struct{ protocol Protocol }

func (a testSignalingAdapter) Name() string       { return string(a.protocol) }
func (a testSignalingAdapter) Protocol() Protocol { return a.protocol }
func (a testSignalingAdapter) MapInvite(context.Context, interface{}) (CallEvent, error) {
	return CallEvent{}, nil
}
func (a testSignalingAdapter) MapAnswer(context.Context, interface{}) (CallEvent, error) {
	return CallEvent{}, nil
}
func (a testSignalingAdapter) MapHold(context.Context, interface{}) (CallEvent, error) {
	return CallEvent{}, nil
}
func (a testSignalingAdapter) MapResume(context.Context, interface{}) (CallEvent, error) {
	return CallEvent{}, nil
}
func (a testSignalingAdapter) MapHangup(context.Context, interface{}) (CallEvent, error) {
	return CallEvent{}, nil
}
func (a testSignalingAdapter) MapTransfer(context.Context, interface{}) (CallEvent, error) {
	return CallEvent{}, nil
}
func (a testSignalingAdapter) MapDTMF(context.Context, interface{}) (CallEvent, error) {
	return CallEvent{}, nil
}

type testMediaAdapter struct{ protocol Protocol }

func (a testMediaAdapter) Name() string       { return string(a.protocol) }
func (a testMediaAdapter) Protocol() Protocol { return a.protocol }
func (a testMediaAdapter) MapSessionStarted(context.Context, interface{}) (MediaEvent, error) {
	return MediaEvent{}, nil
}
func (a testMediaAdapter) MapSessionClosed(context.Context, interface{}) (MediaEvent, error) {
	return MediaEvent{}, nil
}
func (a testMediaAdapter) MapTrackMuted(context.Context, interface{}) (MediaEvent, error) {
	return MediaEvent{}, nil
}
func (a testMediaAdapter) MapTrackUnmuted(context.Context, interface{}) (MediaEvent, error) {
	return MediaEvent{}, nil
}
func (a testMediaAdapter) MapRecordingStarted(context.Context, interface{}) (MediaEvent, error) {
	return MediaEvent{}, nil
}
func (a testMediaAdapter) MapRecordingStopped(context.Context, interface{}) (MediaEvent, error) {
	return MediaEvent{}, nil
}

func TestRegistryRegistersProtocols(t *testing.T) {
	registry := NewRegistry()
	if err := registry.RegisterSignaling(testSignalingAdapter{protocol: ProtocolSIP}); err != nil {
		t.Fatalf("register signaling: %v", err)
	}
	if err := registry.RegisterMedia(testMediaAdapter{protocol: ProtocolWebRTC}); err != nil {
		t.Fatalf("register media: %v", err)
	}

	if _, ok := registry.Signaling(ProtocolSIP); !ok {
		t.Fatalf("expected SIP signaling adapter")
	}
	if _, ok := registry.Media(ProtocolWebRTC); !ok {
		t.Fatalf("expected WebRTC media adapter")
	}
	protocols := registry.SupportedProtocols()
	if len(protocols) != 2 {
		t.Fatalf("expected two protocols, got %v", protocols)
	}
}

func TestRegistryRejectsNilAdapter(t *testing.T) {
	registry := NewRegistry()
	if err := registry.RegisterSignaling(nil); err == nil {
		t.Fatalf("expected nil signaling adapter error")
	}
	if err := registry.RegisterMedia(nil); err == nil {
		t.Fatalf("expected nil media adapter error")
	}
}
