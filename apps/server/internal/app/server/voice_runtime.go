package server

import (
	"fmt"
	"strings"

	"servify/apps/server/internal/config"
	voiceapp "servify/apps/server/internal/modules/voice/application"
	voiceproviderdisabled "servify/apps/server/internal/modules/voice/provider/disabled"
	voiceprovidermock "servify/apps/server/internal/modules/voice/provider/mock"

	"github.com/sirupsen/logrus"
)

const (
	voiceProviderDisabled = "disabled"
	voiceProviderMock     = "mock"
)

func buildVoiceRecordingProvider(cfg *config.Config, logger *logrus.Logger) (voiceapp.RecordingProvider, error) {
	if cfg == nil {
		cfg = config.GetDefaultConfig()
	}
	if logger == nil {
		logger = logrus.StandardLogger()
	}

	provider := normalizeVoiceProvider(cfg.Voice.RecordingProvider)
	switch provider {
	case voiceProviderDisabled:
		logger.Warn("voice recording provider is disabled; recording endpoints will return service unavailable until an explicit provider is configured")
		return voiceproviderdisabled.NewRecordingProvider(), nil
	case voiceProviderMock:
		if strings.EqualFold(strings.TrimSpace(cfg.Server.Environment), "production") {
			return nil, fmt.Errorf("voice recording provider %q is not supported in production", cfg.Voice.RecordingProvider)
		}
		logger.Warn("voice recording provider is using mock implementation; only suitable for development and tests")
		return voiceprovidermock.NewRecordingProvider(), nil
	default:
		return nil, fmt.Errorf("unsupported voice recording provider %q", cfg.Voice.RecordingProvider)
	}
}

func buildVoiceTranscriptProvider(cfg *config.Config, logger *logrus.Logger) (voiceapp.TranscriptProvider, error) {
	if cfg == nil {
		cfg = config.GetDefaultConfig()
	}
	if logger == nil {
		logger = logrus.StandardLogger()
	}

	provider := normalizeVoiceProvider(cfg.Voice.TranscriptProvider)
	switch provider {
	case voiceProviderDisabled:
		logger.Warn("voice transcript provider is disabled; transcript endpoints will return service unavailable until an explicit provider is configured")
		return voiceproviderdisabled.NewTranscriptProvider(), nil
	case voiceProviderMock:
		if strings.EqualFold(strings.TrimSpace(cfg.Server.Environment), "production") {
			return nil, fmt.Errorf("voice transcript provider %q is not supported in production", cfg.Voice.TranscriptProvider)
		}
		logger.Warn("voice transcript provider is using mock implementation; only suitable for development and tests")
		return voiceprovidermock.NewTranscriptProvider(), nil
	default:
		return nil, fmt.Errorf("unsupported voice transcript provider %q", cfg.Voice.TranscriptProvider)
	}
}

func normalizeVoiceProvider(provider string) string {
	provider = strings.TrimSpace(strings.ToLower(provider))
	if provider == "" {
		return voiceProviderDisabled
	}
	return provider
}
