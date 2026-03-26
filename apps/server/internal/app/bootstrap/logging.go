package bootstrap

import (
	"servify/apps/server/internal/config"

	"github.com/sirupsen/logrus"
)

// InitLogging initializes the standard logger from config and returns it.
func InitLogging(cfg *config.Config) (*logrus.Logger, error) {
	if cfg == nil {
		cfg = config.GetDefaultConfig()
	}
	if err := config.InitLogger(cfg); err != nil {
		return nil, err
	}
	logger := logrus.StandardLogger()
	LogSecurityWarnings(logger, cfg)
	return logger, nil
}
