package delivery

import (
	"context"
	"fmt"

	voiceapp "servify/apps/server/internal/modules/voice/application"
	"servify/apps/server/internal/platform/voiceprotocol"
)

// Coordinator is the runtime facade for voice call/media orchestration.
type Coordinator struct {
	calls       *voiceapp.Service
	recordings  *voiceapp.RecordingService
	transcripts *voiceapp.TranscriptService
}

func NewCoordinator(calls *voiceapp.Service, recordings *voiceapp.RecordingService, transcripts *voiceapp.TranscriptService) *Coordinator {
	return &Coordinator{
		calls:       calls,
		recordings:  recordings,
		transcripts: transcripts,
	}
}

func (c *Coordinator) StartCall(ctx context.Context, sessionID string, connectionID string) {
	if c == nil || c.calls == nil {
		return
	}
	_, _ = c.calls.StartCall(ctx, voiceapp.StartCallCommand{
		CallID:       connectionID,
		SessionID:    sessionID,
		ConnectionID: connectionID,
	})
}

func (c *Coordinator) HandleCallEvent(ctx context.Context, event voiceprotocol.CallEvent) error {
	if c == nil || c.calls == nil {
		return fmt.Errorf("voice coordinator unavailable")
	}
	switch event.Kind {
	case voiceprotocol.CallEventInvite:
		_, err := c.calls.StartCall(ctx, voiceapp.StartCallCommand{
			CallID:       event.CallID,
			SessionID:    firstNonEmpty(event.ConversationID, event.ConnectionID, event.CallID),
			ConnectionID: event.ConnectionID,
		})
		return err
	case voiceprotocol.CallEventAnswer:
		_, err := c.calls.AnswerCall(ctx, voiceapp.AnswerCallCommand{CallID: event.CallID})
		return err
	case voiceprotocol.CallEventHold:
		_, err := c.calls.HoldCall(ctx, voiceapp.HoldCallCommand{CallID: event.CallID})
		return err
	case voiceprotocol.CallEventResume:
		_, err := c.calls.ResumeCall(ctx, voiceapp.ResumeCallCommand{CallID: event.CallID})
		return err
	case voiceprotocol.CallEventTransfer:
		return nil
	case voiceprotocol.CallEventHangup:
		_, err := c.calls.EndCall(ctx, voiceapp.EndCallCommand{CallID: event.CallID})
		return err
	case voiceprotocol.CallEventDTMF:
		return nil
	default:
		return fmt.Errorf("unsupported call event kind %q", event.Kind)
	}
}

func (c *Coordinator) HandleMediaEvent(ctx context.Context, event voiceprotocol.MediaEvent) error {
	if c == nil {
		return fmt.Errorf("voice coordinator unavailable")
	}
	switch event.Kind {
	case voiceprotocol.MediaEventRecordingStart:
		if c.recordings == nil {
			return nil
		}
		_, err := c.recordings.StartRecording(ctx, voiceapp.StartRecordingCommand{
			CallID:   event.CallID,
			Provider: string(event.Protocol),
		})
		return err
	case voiceprotocol.MediaEventRecordingStop:
		if c.recordings == nil {
			return nil
		}
		recordingID, _ := event.Metadata["recording_id"].(string)
		if recordingID == "" {
			return nil
		}
		return c.recordings.StopRecording(ctx, voiceapp.StopRecordingCommand{RecordingID: recordingID})
	default:
		return nil
	}
}

func (c *Coordinator) AnswerCall(ctx context.Context, connectionID string) {
	if c == nil || c.calls == nil {
		return
	}
	_, _ = c.calls.AnswerCall(ctx, voiceapp.AnswerCallCommand{CallID: connectionID})
}

func (c *Coordinator) EndCall(ctx context.Context, connectionID string) {
	if c == nil || c.calls == nil {
		return
	}
	_, _ = c.calls.EndCall(ctx, voiceapp.EndCallCommand{CallID: connectionID})
}

func (c *Coordinator) StartRecording(ctx context.Context, cmd voiceapp.StartRecordingCommand) (*voiceapp.RecordingDTO, error) {
	if c == nil || c.recordings == nil {
		return nil, nil
	}
	return c.recordings.StartRecording(ctx, cmd)
}

func (c *Coordinator) StopRecording(ctx context.Context, cmd voiceapp.StopRecordingCommand) error {
	if c == nil || c.recordings == nil {
		return nil
	}
	return c.recordings.StopRecording(ctx, cmd)
}

func (c *Coordinator) AppendTranscript(ctx context.Context, cmd voiceapp.AppendTranscriptCommand) (*voiceapp.TranscriptDTO, error) {
	if c == nil || c.transcripts == nil {
		return nil, nil
	}
	return c.transcripts.Append(ctx, cmd)
}

func (c *Coordinator) GetRecording(ctx context.Context, recordingID string) (*voiceapp.RecordingDTO, error) {
	if c == nil || c.recordings == nil {
		return nil, nil
	}
	return c.recordings.GetRecording(ctx, recordingID)
}

func (c *Coordinator) ListTranscripts(ctx context.Context, callID string) ([]voiceapp.TranscriptDTO, error) {
	if c == nil || c.transcripts == nil {
		return nil, nil
	}
	return c.transcripts.ListByCallID(ctx, callID)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
