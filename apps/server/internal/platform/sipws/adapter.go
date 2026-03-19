package sipws

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"servify/apps/server/internal/platform/voiceprotocol"
)

type Adapter struct{}

func NewAdapter() *Adapter {
	return &Adapter{}
}

func (a *Adapter) Name() string {
	return "sip-ws"
}

func (a *Adapter) Protocol() voiceprotocol.Protocol {
	return voiceprotocol.ProtocolSIPWebSocket
}

func (a *Adapter) MapInvite(_ context.Context, payload interface{}) (voiceprotocol.CallEvent, error) {
	msg, err := asMessage(payload)
	if err != nil {
		return voiceprotocol.CallEvent{}, err
	}
	return toCallEvent(msg, voiceprotocol.CallEventInvite), nil
}

func (a *Adapter) MapAnswer(_ context.Context, payload interface{}) (voiceprotocol.CallEvent, error) {
	msg, err := asMessage(payload)
	if err != nil {
		return voiceprotocol.CallEvent{}, err
	}
	return toCallEvent(msg, voiceprotocol.CallEventAnswer), nil
}

func (a *Adapter) MapHangup(_ context.Context, payload interface{}) (voiceprotocol.CallEvent, error) {
	msg, err := asMessage(payload)
	if err != nil {
		return voiceprotocol.CallEvent{}, err
	}
	return toCallEvent(msg, voiceprotocol.CallEventHangup), nil
}

func (a *Adapter) MapTransfer(_ context.Context, payload interface{}) (voiceprotocol.CallEvent, error) {
	msg, err := asMessage(payload)
	if err != nil {
		return voiceprotocol.CallEvent{}, err
	}
	return toCallEvent(msg, voiceprotocol.CallEventTransfer), nil
}

func (a *Adapter) MapDTMF(_ context.Context, payload interface{}) (voiceprotocol.CallEvent, error) {
	msg, err := asMessage(payload)
	if err != nil {
		return voiceprotocol.CallEvent{}, err
	}
	event := toCallEvent(msg, voiceprotocol.CallEventDTMF)
	event.Metadata["digits"] = msg.DTMF
	return event, nil
}

func asMessage(payload interface{}) (SignalingMessage, error) {
	switch v := payload.(type) {
	case SignalingMessage:
		return v, nil
	case map[string]interface{}:
		raw, err := json.Marshal(v)
		if err != nil {
			return SignalingMessage{}, fmt.Errorf("marshal sip-ws payload: %w", err)
		}
		var msg SignalingMessage
		if err := json.Unmarshal(raw, &msg); err != nil {
			return SignalingMessage{}, fmt.Errorf("decode sip-ws payload: %w", err)
		}
		return msg, nil
	default:
		return SignalingMessage{}, fmt.Errorf("unsupported sip-ws payload type %T", payload)
	}
}

func toCallEvent(msg SignalingMessage, kind voiceprotocol.CallEventKind) voiceprotocol.CallEvent {
	occurredAt := msg.OccurredAt
	if occurredAt.IsZero() {
		occurredAt = time.Now()
	}
	metadata := map[string]interface{}{}
	for k, v := range msg.Metadata {
		metadata[k] = v
	}
	if msg.Method != "" {
		metadata["method"] = msg.Method
	}
	return voiceprotocol.CallEvent{
		EventID:        fmt.Sprintf("sipws-%s-%s", kind, msg.CallID),
		Protocol:       voiceprotocol.ProtocolSIPWebSocket,
		Kind:           kind,
		CallID:         msg.CallID,
		ConversationID: msg.ConversationID,
		ConnectionID:   msg.ConnectionID,
		From:           msg.From,
		To:             msg.To,
		OccurredAt:     occurredAt,
		Metadata:       metadata,
	}
}

var _ voiceprotocol.CallSignalingAdapter = (*Adapter)(nil)
