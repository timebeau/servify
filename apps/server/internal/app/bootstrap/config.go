package bootstrap

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"servify/apps/server/internal/config"

	"github.com/spf13/viper"
)

type RuntimeOverrides struct {
	Database DatabaseOptions
	HTTP     HTTPServerOptions
}

// LoadConfig loads configuration from the default path or a specific config file.
func LoadConfig(configPath string) (*config.Config, error) {
	viper.Reset()
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	viper.AddConfigPath(filepath.Join("..", ".."))
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
	cfg, result, err := config.LoadWithResult()
	if err != nil {
		return nil, err
	}
	// Log security warnings for development environments
	if len(result.Warnings) > 0 && cfg.Server.Environment != "production" {
		LogSecurityWarnings(nil, cfg)
	}
	return cfg, nil
}

// ResolveRuntimeOverrides parses server startup flags and environment overrides.
func ResolveRuntimeOverrides(cfg *config.Config, args []string, output io.Writer) (RuntimeOverrides, error) {
	if cfg == nil {
		cfg = config.GetDefaultConfig()
	}

	var (
		dbDriver  string
		flagDSN   string
		dbHost    string
		dbPortStr string
		dbUser    string
		dbPass    string
		dbName    string
		dbSSLMode string
		dbTZ      string
		srvHost   string
		srvPort   int
	)

	flagSet := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	if output != nil {
		flagSet.SetOutput(output)
	} else {
		flagSet.SetOutput(io.Discard)
	}
	flagSet.StringVar(&dbDriver, "db-driver", getenvDefault("DB_DRIVER", "postgres"), "database driver (postgres or sqlite)")
	defaultDSN := ""
	if strings.EqualFold(getenvDefault("DB_DRIVER", "postgres"), "sqlite") {
		defaultDSN = os.Getenv("DB_DSN")
	}
	flagSet.StringVar(&flagDSN, "dsn", defaultDSN, "database DSN; when db-driver=sqlite this should be a sqlite file path or DSN")
	flagSet.StringVar(&dbHost, "db-host", getenvDefault("DB_HOST", cfg.Database.Host), "database host")
	flagSet.StringVar(&dbPortStr, "db-port", getenvDefault("DB_PORT", fmt.Sprintf("%d", cfg.Database.Port)), "database port")
	flagSet.StringVar(&dbUser, "db-user", getenvDefault("DB_USER", cfg.Database.User), "database user")
	flagSet.StringVar(&dbPass, "db-pass", getenvDefault("DB_PASSWORD", cfg.Database.Password), "database password")
	flagSet.StringVar(&dbName, "db-name", getenvDefault("DB_NAME", cfg.Database.Name), "database name")
	flagSet.StringVar(&dbSSLMode, "db-sslmode", getenvDefault("DB_SSLMODE", "disable"), "sslmode (disable, require, verify-ca, verify-full)")
	flagSet.StringVar(&dbTZ, "db-timezone", getenvDefault("DB_TIMEZONE", "UTC"), "database timezone")
	flagSet.StringVar(&srvHost, "host", getenvDefault("SERVIFY_HOST", cfg.Server.Host), "server host (listen)")
	flagSet.IntVar(&srvPort, "port", serverPortFromEnv(cfg.Server.Port), "server port (listen)")
	if err := flagSet.Parse(args); err != nil {
		return RuntimeOverrides{}, err
	}

	dbOpts := DatabaseOptions{
		Driver:   dbDriver,
		DSN:      flagDSN,
		Host:     dbHost,
		Port:     dbPortStr,
		User:     dbUser,
		Password: dbPass,
		Name:     dbName,
		SSLMode:  dbSSLMode,
		TimeZone: dbTZ,
	}
	if dbOpts.DSN == "" && normalizedDatabaseDriver(dbOpts) != "sqlite" {
		dbOpts.DSN = BuildPostgresDSN(cfg, dbOpts)
	}

	return RuntimeOverrides{
		Database: dbOpts,
		HTTP: HTTPServerOptions{
			Host: srvHost,
			Port: srvPort,
		},
	}, nil
}

func getenvDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func serverPortFromEnv(defaultPort int) int {
	if p := os.Getenv("SERVIFY_PORT"); p != "" {
		if n, err := strconv.Atoi(p); err == nil {
			return n
		}
	}
	return defaultPort
}
