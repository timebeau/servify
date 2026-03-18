package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	appbootstrap "servify/apps/server/internal/app/bootstrap"
	appserver "servify/apps/server/internal/app/server"
	appworker "servify/apps/server/internal/app/worker"
	"servify/apps/server/internal/platform/eventbus"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
)

func main() {
	cfg, err := appbootstrap.LoadConfig("")
	if err != nil {
		logrus.Fatalf("Failed to load config: %v", err)
	}
	app := appbootstrap.BuildApp(cfg)

	// 允许通过 flags/env 覆盖数据库连接（保持与 migrate 一致的接口）
	var (
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
	// 延迟导入以避免顶层 import 冲突
	{
		// 标准库 flag 在此作用域使用
		type strptr = *string
		_ = strptr(nil)
	}
	// 使用标准库 flag
	flagSet := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flagSet.SetOutput(os.Stdout)
	flagSet.StringVar(&flagDSN, "dsn", os.Getenv("DB_DSN"), "Postgres DSN, if set overrides other DB flags")
	flagSet.StringVar(&dbHost, "db-host", getenvDefault("DB_HOST", cfg.Database.Host), "database host")
	flagSet.StringVar(&dbPortStr, "db-port", getenvDefault("DB_PORT", fmt.Sprintf("%d", cfg.Database.Port)), "database port")
	flagSet.StringVar(&dbUser, "db-user", getenvDefault("DB_USER", cfg.Database.User), "database user")
	flagSet.StringVar(&dbPass, "db-pass", getenvDefault("DB_PASSWORD", cfg.Database.Password), "database password")
	flagSet.StringVar(&dbName, "db-name", getenvDefault("DB_NAME", cfg.Database.Name), "database name")
	flagSet.StringVar(&dbSSLMode, "db-sslmode", getenvDefault("DB_SSLMODE", "disable"), "sslmode (disable, require, verify-ca, verify-full)")
	flagSet.StringVar(&dbTZ, "db-timezone", getenvDefault("DB_TIMEZONE", "UTC"), "database timezone")
	flagSet.StringVar(&srvHost, "host", getenvDefault("SERVIFY_HOST", cfg.Server.Host), "server host (listen)")
	flagSet.IntVar(&srvPort, "port", func() int {
		if p := os.Getenv("SERVIFY_PORT"); p != "" {
			if n, err := strconv.Atoi(p); err == nil {
				return n
			}
		}
		return cfg.Server.Port
	}(), "server port (listen)")
	_ = flagSet.Parse(os.Args[1:])

	// 组装 DSN
	dsn := flagDSN
	if dsn == "" {
		host := firstNonEmpty(dbHost, cfg.Database.Host)
		user := firstNonEmpty(dbUser, cfg.Database.User)
		pass := firstNonEmpty(dbPass, cfg.Database.Password)
		name := firstNonEmpty(dbName, cfg.Database.Name)
		port := dbPortStr
		if port == "" && cfg.Database.Port != 0 {
			port = fmt.Sprintf("%d", cfg.Database.Port)
		}
		ssl := dbSSLMode
		tz := dbTZ
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", host, user, pass, name, port, ssl, tz)
	}
	appLogger, err := appbootstrap.InitLogging(cfg)
	if err != nil {
		logrus.Warnf("init logger: %v", err)
	}
	app.Logger = appLogger

	if err := appbootstrap.SetupObservability(context.Background(), cfg, app); err != nil {
		appLogger.Warnf("init tracing: %v", err)
	}

	db, err := appbootstrap.OpenDatabase(cfg, appbootstrap.DatabaseOptions{
		DSN:           dsn,
		LogLevel:      logger.Info,
		EnableTracing: cfg.Monitoring.Tracing.Enabled,
	})
	if err != nil {
		appLogger.Fatalf("Failed to connect to database: %v", err)
	}
	app.DB = db

	if appbootstrap.AutoMigrateEnabled() {
		if err := appbootstrap.AutoMigrate(db); err != nil {
			appLogger.Fatalf("Failed to migrate database: %v", err)
		}
	}

	bus := eventbus.NewInMemoryBus()
	app.EventBus = bus

	runtime, err := appserver.BuildRuntime(cfg, appLogger, db, bus)
	if err != nil {
		appLogger.Fatalf("Failed to build runtime: %v", err)
	}
	if err := runtime.Start(); err != nil {
		appLogger.Fatalf("Failed to start runtime: %v", err)
	}

	app.RegisterWorker(appworker.NewStatisticsWorker(runtime.StatisticsService, time.Hour))
	app.RegisterWorker(appworker.NewSLAMonitorWorker(runtime.SLAService, 5*time.Minute))
	if err := app.StartWorkers(); err != nil {
		appLogger.Fatalf("Failed to start workers: %v", err)
	}

	r := appserver.BuildRouter(runtime.RouterDependencies())
	app.Router = r

	srv := appbootstrap.NewHTTPServer(cfg, r, appbootstrap.HTTPServerOptions{
		Host: firstNonEmpty(srvHost, cfg.Server.Host),
		Port: srvPort,
	})
	appbootstrap.StartHTTPServer(srv, appLogger, fmt.Sprintf("Starting server on %s", srv.Addr))

	appbootstrap.WaitForShutdownSignal()
	appLogger.Info("Shutting down server...")
	shutdownCtx, cancel := appbootstrap.ShutdownContext(30 * time.Second)
	defer cancel()
	if err := app.StopWorkers(shutdownCtx); err != nil {
		appLogger.Errorf("Failed to stop workers cleanly: %v", err)
	}
	if err := runtime.Stop(shutdownCtx); err != nil {
		appLogger.Errorf("Failed to stop runtime cleanly: %v", err)
	}
	if err := app.RunShutdownHooks(); err != nil {
		appLogger.Errorf("Failed to run shutdown hooks: %v", err)
	}
	if err := srv.Shutdown(shutdownCtx); err != nil {
		appLogger.Fatalf("Server forced to shutdown: %v", err)
	}
	appLogger.Info("Server exited")
}

// helpers (copied from migrate for consistency)
func getenvDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}
