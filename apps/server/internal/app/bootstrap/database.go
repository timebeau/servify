package bootstrap

import (
	"fmt"

	"servify/apps/server/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	gormtracing "gorm.io/plugin/opentelemetry/tracing"
)

type DatabaseOptions struct {
	DSN           string
	Host          string
	Port          string
	User          string
	Password      string
	Name          string
	SSLMode       string
	TimeZone      string
	LogLevel      logger.LogLevel
	EnableTracing bool
}

// BuildPostgresDSN composes a DSN from explicit options and config defaults.
func BuildPostgresDSN(cfg *config.Config, opts DatabaseOptions) string {
	if opts.DSN != "" {
		return opts.DSN
	}
	dbCfg := config.DatabaseConfig{}
	if cfg != nil {
		dbCfg = cfg.Database
	}
	host := firstNonEmpty(opts.Host, dbCfg.Host)
	user := firstNonEmpty(opts.User, dbCfg.User)
	pass := firstNonEmpty(opts.Password, dbCfg.Password)
	name := firstNonEmpty(opts.Name, dbCfg.Name)
	port := firstNonEmpty(opts.Port, fmt.Sprintf("%d", dbCfg.Port))
	sslMode := firstNonEmpty(opts.SSLMode, "disable")
	timeZone := firstNonEmpty(opts.TimeZone, "UTC")
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		host, user, pass, name, port, sslMode, timeZone,
	)
}

// OpenDatabase opens GORM and injects tracing when requested.
func OpenDatabase(cfg *config.Config, opts DatabaseOptions) (*gorm.DB, error) {
	level := opts.LogLevel
	if level == 0 {
		level = logger.Warn
	}
	db, err := gorm.Open(postgres.Open(BuildPostgresDSN(cfg, opts)), &gorm.Config{
		Logger:                                   logger.Default.LogMode(level),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, err
	}
	enableTracing := opts.EnableTracing
	if cfg != nil && cfg.Monitoring.Tracing.Enabled {
		enableTracing = true
	}
	if enableTracing {
		_ = db.Use(gormtracing.NewPlugin())
	}
	return db, nil
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}
