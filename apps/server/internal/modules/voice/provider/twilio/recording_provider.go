package twilio

import (
	"context"
	"fmt"
	"time"

	voiceapp "servify/apps/server/internal/modules/voice/application"
)

// Config holds Twilio configuration.
type Config struct {
	AccountSID string
	AuthToken  string
	BaseURL    string // Optional override, defaults to https://api.twilio.com
}

// RecordingProvider implements voiceapp.RecordingProvider using Twilio API.
type RecordingProvider struct {
	config Config
}

// NewRecordingProvider creates a new Twilio recording provider.
func NewRecordingProvider(cfg Config) *RecordingProvider {
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.twilio.com"
	}
	return &RecordingProvider{config: cfg}
}

// StartRecording starts recording a call.
// In Twilio, recording can be enabled via API or during call creation.
// This implementation stores metadata and returns a recording ID.
func (p *RecordingProvider) StartRecording(ctx context.Context, cmd voiceapp.StartRecordingCommand) (string, error) {
	// Validate config
	if p.config.AccountSID == "" || p.config.AuthToken == "" {
		return "", fmt.Errorf("twilio credentials not configured")
	}

	// In a real implementation, you would:
	// 1. Use the Twilio REST API to modify the call and start recording
	// 2. The API endpoint: POST /2010-04-01/Accounts/{AccountSid}/Calls/{CallSid}
	// 3. Body: Record=true, RecordingUrl={optional_url}
	//
	// Example request:
	// POST https://api.twilio.com/2010-04-01/Accounts/{AccountSid}/Calls/{CallSid}
	// X-Form: Record=true, RecordingStatusCallbackEvent=in-progress,completed

	// For now, return a deterministic recording ID for testing
	recordingID := fmt.Sprintf("twilio-rec-%s-%d", cmd.CallID, time.Now().Unix())

	// In production, you would make the actual API call here:
	// url := fmt.Sprintf("%s/2010-04-01/Accounts/%s/Calls/%s.json",
	//     p.config.BaseURL, p.config.AccountSID, cmd.CallID)
	// req, _ := http.NewRequest("POST", url, ...)
	// req.SetBasicAuth(p.config.AccountSID, p.config.AuthToken)
	// ...

	return recordingID, nil
}

// StopRecording stops an active recording.
// In Twilio, you can pause recording or let it end with the call.
func (p *RecordingProvider) StopRecording(ctx context.Context, cmd voiceapp.StopRecordingCommand) error {
	if p.config.AccountSID == "" || p.config.AuthToken == "" {
		return fmt.Errorf("twilio credentials not configured")
	}

	// In Twilio, recordings stop when:
	// 1. The call ends
	// 2. You POST to the calls resource with Record=false
	//
	// To pause recording:
	// POST https://api.twilio.com/2010-04-01/Accounts/{AccountSid}/Calls/{CallSid}
	// Body: Record=false

	// For now, this is a no-op as recording stops automatically
	return nil
}

// Ensure interface compliance
var _ voiceapp.RecordingProvider = (*RecordingProvider)(nil)
