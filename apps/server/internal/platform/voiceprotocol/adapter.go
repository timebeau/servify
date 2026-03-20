package voiceprotocol

import "context"

// CallSignalingAdapter maps protocol-specific signaling into normalized call events.
type CallSignalingAdapter interface {
	Name() string
	Protocol() Protocol
	MapInvite(context.Context, interface{}) (CallEvent, error)
	MapAnswer(context.Context, interface{}) (CallEvent, error)
	MapHold(context.Context, interface{}) (CallEvent, error)
	MapResume(context.Context, interface{}) (CallEvent, error)
	MapHangup(context.Context, interface{}) (CallEvent, error)
	MapTransfer(context.Context, interface{}) (CallEvent, error)
	MapDTMF(context.Context, interface{}) (CallEvent, error)
}

// MediaSessionAdapter maps media-plane lifecycle into normalized media events.
type MediaSessionAdapter interface {
	Name() string
	Protocol() Protocol
	MapSessionStarted(context.Context, interface{}) (MediaEvent, error)
	MapSessionClosed(context.Context, interface{}) (MediaEvent, error)
	MapTrackMuted(context.Context, interface{}) (MediaEvent, error)
	MapTrackUnmuted(context.Context, interface{}) (MediaEvent, error)
	MapRecordingStarted(context.Context, interface{}) (MediaEvent, error)
	MapRecordingStopped(context.Context, interface{}) (MediaEvent, error)
}
