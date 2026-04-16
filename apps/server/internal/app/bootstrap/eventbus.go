package bootstrap

import (
	"fmt"
	"strings"

	"servify/apps/server/internal/config"
	"servify/apps/server/internal/platform/eventbus"

	"github.com/sirupsen/logrus"
)

const eventBusProviderInMemory = "inmemory"

func BuildEventBus(cfg *config.Config, logger *logrus.Logger) (eventbus.Bus, error) {
	if cfg == nil {
		cfg = config.GetDefaultConfig()
	}
	if logger == nil {
		logger = logrus.StandardLogger()
	}

	provider := strings.TrimSpace(strings.ToLower(cfg.EventBus.Provider))
	if provider == "" {
		provider = eventBusProviderInMemory
	}

	switch provider {
	case eventBusProviderInMemory:
		if strings.EqualFold(strings.TrimSpace(cfg.Server.Environment), "production") {
			logger.Warn("event bus provider 'inmemory' is running in production; asynchronous events are not durable and in-flight events are lost on restart")
		}
		return eventbus.NewInMemoryBusWithLogger(logger), nil
	default:
		return nil, fmt.Errorf("unsupported event bus provider %q", cfg.EventBus.Provider)
	}
}
