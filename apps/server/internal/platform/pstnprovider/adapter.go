package pstnprovider

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
	return "pstn-provider"
}

func (a *Adapter) Protocol() voiceprotocol.Protocol {
	return voiceprotocol.ProtocolPSTNProvider
}

func (a *Adapter) MapInvite(_ context.Context, payload interface{}) (voiceprotocol.CallEvent, error) {
	event, err := asWebhook(payload)
	if err != nil {
		return voiceprotocol.CallEvent{}, err
	}
	return toCallEvent(event, voiceprotocol.CallEventInvite), nil
}

func (a *Adapter) MapAnswer(_ context.Context, payload interface{}) (voiceprotocol.CallEvent, error) {
	event, err := asWebhook(payload)
	if err != nil {
		return voiceprotocol.CallEvent{}, err
	}
	return toCallEvent(event, voiceprotocol.CallEventAnswer), nil
}

func (a *Adapter) MapHangup(_ context.Context, payload interface{}) (voiceprotocol.CallEvent, error) {
	event, err := asWebhook(payload)
	if err != nil {
		return voiceprotocol.CallEvent{}, err
	}
	return toCallEvent(event, voiceprotocol.CallEventHangup), nil
}

func (a *Adapter) MapTransfer(_ context.Context, payload interface{}) (voiceprotocol.CallEvent, error) {
	event, err := asWebhook(payload)
	if err != nil {
		return voiceprotocol.CallEvent{}, err
	}
	return toCallEvent(event, voiceprotocol.CallEventTransfer), nil
}

func (a *Adapter) MapDTMF(_ context.Context, payload interface{}) (voiceprotocol.CallEvent, error) {
	event, err := asWebhook(payload)
	if err != nil {
		return voiceprotocol.CallEvent{}, err
	}
	callEvent := toCallEvent(event, voiceprotocol.CallEventDTMF)
	callEvent.Metadata["digits"] = event.DTMF
	return callEvent, nil
}

func asWebhook(payload interface{}) (WebhookEvent, error) {
	switch v := payload.(type) {
	case WebhookEvent:
		return v, nil
	case map[string]interface{}:
		raw, err := json.Marshal(v)
		if err != nil {
			return WebhookEvent{}, fmt.Errorf("marshal pstn payload: %w", err)
		}
		var event WebhookEvent
		if err := json.Unmarshal(raw, &event); err != nil {
			return WebhookEvent{}, fmt.Errorf("decode pstn payload: %w", err)
		}
		return event, nil
	default:
		return WebhookEvent{}, fmt.Errorf("unsupported pstn payload type %T", payload)
	}
}

func toCallEvent(event WebhookEvent, kind voiceprotocol.CallEventKind) voiceprotocol.CallEvent {
	occurredAt := event.OccurredAt
	if occurredAt.IsZero() {
		occurredAt = time.Now()
	}
	metadata := map[string]interface{}{}
	for k, v := range event.Metadata {
		metadata[k] = v
	}
	if event.Provider != "" {
		metadata["provider"] = event.Provider
	}
	if event.EventType != "" {
		metadata["event_type"] = event.EventType
	}
	eventID := event.EventID
	if eventID == "" {
		eventID = fmt.Sprintf("pstn-%s-%s", kind, event.CallID)
	}
	return voiceprotocol.CallEvent{
		EventID:        eventID,
		Protocol:       voiceprotocol.ProtocolPSTNProvider,
		Kind:           kind,
		CallID:         event.CallID,
		ConversationID: event.ConversationID,
		From:           event.From,
		To:             event.To,
		OccurredAt:     occurredAt,
		Metadata:       metadata,
	}
}

var _ voiceprotocol.CallSignalingAdapter = (*Adapter)(nil)
