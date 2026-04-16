package server

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"servify/apps/server/internal/config"
	voiceapp "servify/apps/server/internal/modules/voice/application"

	"github.com/sirupsen/logrus"
)

func TestBuildVoiceRecordingProviderDefaultsToDisabled(t *testing.T) {
	cfg := config.GetDefaultConfig()
	logger := newVoiceRuntimeTestLogger(nil)

	provider, err := buildVoiceRecordingProvider(cfg, logger)
	if err != nil {
		t.Fatalf("buildVoiceRecordingProvider() error = %v", err)
	}
	if provider == nil {
		t.Fatal("expected recording provider")
	}
	if _, err := provider.StartRecording(context.Background(), voiceapp.StartRecordingCommand{CallID: "call-1"}); err == nil {
		t.Fatal("expected disabled provider error")
	}
}

func TestBuildVoiceProvidersRejectMockInProduction(t *testing.T) {
	cfg := config.GetDefaultConfig()
	cfg.Server.Environment = "production"
	cfg.Voice.RecordingProvider = "mock"
	cfg.Voice.TranscriptProvider = "mock"

	if _, err := buildVoiceRecordingProvider(cfg, newVoiceRuntimeTestLogger(nil)); err == nil {
		t.Fatal("expected production mock recording provider to fail")
	}
	if _, err := buildVoiceTranscriptProvider(cfg, newVoiceRuntimeTestLogger(nil)); err == nil {
		t.Fatal("expected production mock transcript provider to fail")
	}
}

func TestBuildVoiceProvidersAllowMockInDevelopment(t *testing.T) {
	cfg := config.GetDefaultConfig()
	cfg.Voice.RecordingProvider = "mock"
	cfg.Voice.TranscriptProvider = "mock"

	var buf bytes.Buffer
	logger := newVoiceRuntimeTestLogger(&buf)

	recordingProvider, err := buildVoiceRecordingProvider(cfg, logger)
	if err != nil {
		t.Fatalf("buildVoiceRecordingProvider() error = %v", err)
	}
	if recordingProvider == nil {
		t.Fatal("expected recording provider")
	}

	transcriptProvider, err := buildVoiceTranscriptProvider(cfg, logger)
	if err != nil {
		t.Fatalf("buildVoiceTranscriptProvider() error = %v", err)
	}
	if transcriptProvider == nil {
		t.Fatal("expected transcript provider")
	}

	if !strings.Contains(buf.String(), "mock implementation") {
		t.Fatalf("expected mock provider warning, got %q", buf.String())
	}
}

func newVoiceRuntimeTestLogger(buf *bytes.Buffer) *logrus.Logger {
	logger := logrus.New()
	if buf != nil {
		logger.SetOutput(buf)
	}
	logger.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
		DisableQuote:     true,
	})
	return logger
}
