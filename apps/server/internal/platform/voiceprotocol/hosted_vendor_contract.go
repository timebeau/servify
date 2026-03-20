package voiceprotocol

import "context"

// HostedVendorWebhookAdapter reserves the contract for future hosted voice
// providers that deliver call control through signed webhook events.
type HostedVendorWebhookAdapter interface {
	CallSignalingAdapter
	WebhookPath() string
	ValidateSignature(ctx context.Context, headers map[string]string, body []byte) error
}
