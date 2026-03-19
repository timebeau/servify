package sip

import (
	"context"

	channelplatform "servify/apps/server/internal/platform/channel"
)

// SIPAdapter maps SIP signaling into normalized channel events.
type SIPAdapter interface {
	Name() string
	MapInvite(ctx context.Context, call InboundCall) (channelplatform.InboundEvent, error)
	MapHangup(ctx context.Context, call InboundCall) (channelplatform.InboundEvent, error)
	MapDTMF(ctx context.Context, call InboundCall) (channelplatform.InboundEvent, error)
}
