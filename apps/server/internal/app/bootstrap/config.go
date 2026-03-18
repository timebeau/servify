package bootstrap

import (
	"errors"

	"servify/apps/server/internal/config"

	"github.com/spf13/viper"
)

// LoadConfig loads configuration from the default path or a specific config file.
func LoadConfig(configPath string) (*config.Config, error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	if configPath != "" {
		viper.SetConfigFile(configPath)
	}

	if err := viper.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if configPath != "" || !errors.As(err, &notFound) {
			return nil, err
		}
	}
	return config.Load(), nil
}
