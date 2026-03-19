package voiceprotocol

import "time"

type Protocol string

const (
	ProtocolSIP          Protocol = "sip"
	ProtocolSIPWebSocket Protocol = "sip-ws"
	ProtocolWebRTC       Protocol = "webrtc"
	ProtocolRTP          Protocol = "rtp"
	ProtocolSRTP         Protocol = "srtp"
	ProtocolPSTNProvider Protocol = "pstn-provider"
)

type CallEventKind string

const (
	CallEventInvite   CallEventKind = "invite"
	CallEventAnswer   CallEventKind = "answer"
	CallEventHangup   CallEventKind = "hangup"
	CallEventTransfer CallEventKind = "transfer"
	CallEventHold     CallEventKind = "hold"
	CallEventResume   CallEventKind = "resume"
	CallEventDTMF     CallEventKind = "dtmf"
)

type MediaEventKind string

const (
	MediaEventSessionStarted MediaEventKind = "session_started"
	MediaEventSessionClosed  MediaEventKind = "session_closed"
	MediaEventTrackPublished MediaEventKind = "track_published"
	MediaEventTrackMuted     MediaEventKind = "track_muted"
	MediaEventTrackUnmuted   MediaEventKind = "track_unmuted"
	MediaEventRecordingStart MediaEventKind = "recording_started"
	MediaEventRecordingStop  MediaEventKind = "recording_stopped"
)

type CallEvent struct {
	EventID        string                 `json:"event_id"`
	Protocol       Protocol               `json:"protocol"`
	Kind           CallEventKind          `json:"kind"`
	CallID         string                 `json:"call_id"`
	ConversationID string                 `json:"conversation_id,omitempty"`
	ConnectionID   string                 `json:"connection_id,omitempty"`
	From           string                 `json:"from,omitempty"`
	To             string                 `json:"to,omitempty"`
	OccurredAt     time.Time              `json:"occurred_at"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

type MediaEvent struct {
	EventID        string                 `json:"event_id"`
	Protocol       Protocol               `json:"protocol"`
	Kind           MediaEventKind         `json:"kind"`
	CallID         string                 `json:"call_id,omitempty"`
	ConversationID string                 `json:"conversation_id,omitempty"`
	ConnectionID   string                 `json:"connection_id,omitempty"`
	OccurredAt     time.Time              `json:"occurred_at"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}
