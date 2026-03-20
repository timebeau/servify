package delivery

import "context"

type MediaBridgeContract interface {
	BridgeSession(ctx context.Context, req BridgeRequest) (*BridgeSession, error)
	CloseBridge(ctx context.Context, bridgeID string) error
}

type BridgeRequest struct {
	CallID            string
	WebRTCConnection  string
	PacketConnection  string
	PacketProtocol    string
}

type BridgeSession struct {
	ID               string `json:"id"`
	CallID           string `json:"call_id"`
	WebRTCConnection string `json:"webrtc_connection"`
	PacketConnection string `json:"packet_connection"`
	PacketProtocol   string `json:"packet_protocol"`
}

type ConferenceMixer interface {
	CreateConference(ctx context.Context, req ConferenceRequest) (*ConferenceSession, error)
	AddParticipant(ctx context.Context, conferenceID string, participantID string) error
	RemoveParticipant(ctx context.Context, conferenceID string, participantID string) error
}

type ConferenceRequest struct {
	CallID string
}

type ConferenceSession struct {
	ID           string   `json:"id"`
	CallID       string   `json:"call_id"`
	Participants []string `json:"participants"`
}

type VoiceQualitySample struct {
	CallID      string  `json:"call_id"`
	PacketLoss  float64 `json:"packet_loss"`
	JitterMs    float64 `json:"jitter_ms"`
	RTTMs       float64 `json:"rtt_ms"`
	MOS         float64 `json:"mos"`
	Source      string  `json:"source"`
}

type VoiceQualityMetricsSink interface {
	RecordQualitySample(ctx context.Context, sample VoiceQualitySample) error
}
