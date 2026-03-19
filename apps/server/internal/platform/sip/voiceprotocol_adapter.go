package sip

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"servify/apps/server/internal/platform/voiceprotocol"
)

// VoiceProtocolAdapter bridges SIP signaling into protocol-neutral voice events.
type VoiceProtocolAdapter struct{}

func NewVoiceProtocolAdapter() *VoiceProtocolAdapter {
	return &VoiceProtocolAdapter{}
}

func (a *VoiceProtocolAdapter) Name() string {
	return ChannelName
}

func (a *VoiceProtocolAdapter) Protocol() voiceprotocol.Protocol {
	return voiceprotocol.ProtocolSIP
}

func (a *VoiceProtocolAdapter) MapInvite(_ context.Context, payload interface{}) (voiceprotocol.CallEvent, error) {
	call, err := asInboundCall(payload)
	if err != nil {
		return voiceprotocol.CallEvent{}, err
	}
	return toCallEvent(call, voiceprotocol.CallEventInvite), nil
}

func (a *VoiceProtocolAdapter) MapAnswer(_ context.Context, payload interface{}) (voiceprotocol.CallEvent, error) {
	call, err := asInboundCall(payload)
	if err != nil {
		return voiceprotocol.CallEvent{}, err
	}
	return toCallEvent(call, voiceprotocol.CallEventAnswer), nil
}

func (a *VoiceProtocolAdapter) MapHangup(_ context.Context, payload interface{}) (voiceprotocol.CallEvent, error) {
	call, err := asInboundCall(payload)
	if err != nil {
		return voiceprotocol.CallEvent{}, err
	}
	return toCallEvent(call, voiceprotocol.CallEventHangup), nil
}

func (a *VoiceProtocolAdapter) MapTransfer(_ context.Context, payload interface{}) (voiceprotocol.CallEvent, error) {
	call, err := asInboundCall(payload)
	if err != nil {
		return voiceprotocol.CallEvent{}, err
	}
	return toCallEvent(call, voiceprotocol.CallEventTransfer), nil
}

func (a *VoiceProtocolAdapter) MapDTMF(_ context.Context, payload interface{}) (voiceprotocol.CallEvent, error) {
	call, err := asInboundCall(payload)
	if err != nil {
		return voiceprotocol.CallEvent{}, err
	}
	event := toCallEvent(call, voiceprotocol.CallEventDTMF)
	event.Metadata["digits"] = call.DTMF
	return event, nil
}

func asInboundCall(payload interface{}) (InboundCall, error) {
	switch v := payload.(type) {
	case InboundCall:
		return v, nil
	case map[string]interface{}:
		raw, err := json.Marshal(v)
		if err != nil {
			return InboundCall{}, fmt.Errorf("marshal SIP payload: %w", err)
		}
		var call InboundCall
		if err := json.Unmarshal(raw, &call); err != nil {
			return InboundCall{}, fmt.Errorf("decode SIP payload: %w", err)
		}
		return call, nil
	default:
		return InboundCall{}, fmt.Errorf("unsupported SIP payload type %T", payload)
	}
}

func toCallEvent(call InboundCall, kind voiceprotocol.CallEventKind) voiceprotocol.CallEvent {
	occurredAt := call.OccurredAt
	if occurredAt.IsZero() {
		occurredAt = time.Now()
	}
	metadata := map[string]interface{}{}
	for k, v := range call.Metadata {
		metadata[k] = v
	}
	if call.DTMF != "" {
		metadata["dtmf"] = call.DTMF
	}
	return voiceprotocol.CallEvent{
		EventID:        fmt.Sprintf("sip-%s-%s", kind, call.CallID),
		Protocol:       voiceprotocol.ProtocolSIP,
		Kind:           kind,
		CallID:         call.CallID,
		ConversationID: call.ConversationID,
		From:           call.From,
		To:             call.To,
		OccurredAt:     occurredAt,
		Metadata:       metadata,
	}
}

var _ voiceprotocol.CallSignalingAdapter = (*VoiceProtocolAdapter)(nil)
