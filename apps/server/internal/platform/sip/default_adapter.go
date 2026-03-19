package sip

import (
	"context"
	"fmt"
	"time"

	channelplatform "servify/apps/server/internal/platform/channel"
)

const ChannelName = "sip"

// DefaultAdapter maps SIP signaling to normalized channel events.
type DefaultAdapter struct{}

func NewDefaultAdapter() *DefaultAdapter {
	return &DefaultAdapter{}
}

func (a *DefaultAdapter) Name() string {
	return ChannelName
}

func (a *DefaultAdapter) MapInvite(_ context.Context, call InboundCall) (channelplatform.InboundEvent, error) {
	return a.mapCallEvent(call, CallEventInvite), nil
}

func (a *DefaultAdapter) MapHangup(_ context.Context, call InboundCall) (channelplatform.InboundEvent, error) {
	return a.mapCallEvent(call, CallEventHangup), nil
}

func (a *DefaultAdapter) MapDTMF(_ context.Context, call InboundCall) (channelplatform.InboundEvent, error) {
	event := a.mapCallEvent(call, CallEventDTMF)
	event.Payload["digits"] = call.DTMF
	return event, nil
}

func (a *DefaultAdapter) mapCallEvent(call InboundCall, event CallEventKind) channelplatform.InboundEvent {
	occurredAt := call.OccurredAt
	if occurredAt.IsZero() {
		occurredAt = time.Now()
	}

	payload := map[string]interface{}{
		"call_id":    call.CallID,
		"event":      event,
		"from":       call.From,
		"to":         call.To,
		"metadata":   call.Metadata,
		"dtmf":       call.DTMF,
		"protocol":   ChannelName,
		"occurredAt": occurredAt,
	}

	return channelplatform.InboundEvent{
		EventID:        fmt.Sprintf("sip-%s-%s", event, call.CallID),
		Channel:        ChannelName,
		ConversationID: call.ConversationID,
		ActorID:        call.From,
		Kind:           channelplatform.EventKindCall,
		Payload:        payload,
		OccurredAt:     occurredAt,
	}
}

var _ SIPAdapter = (*DefaultAdapter)(nil)
