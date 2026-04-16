package bootstrap

import (
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"servify/apps/server/internal/config"
)

func TestBuildPostgresDSN(t *testing.T) {
	cfg := config.GetDefaultConfig()
	dsn := BuildPostgresDSN(cfg, DatabaseOptions{
		Host:     "db.local",
		Port:     "5433",
		User:     "svc",
		Password: "secret",
		Name:     "servify_test",
		SSLMode:  "require",
		TimeZone: "Asia/Shanghai",
	})
	for _, part := range []string{
		"host=db.local",
		"user=svc",
		"password=secret",
		"dbname=servify_test",
		"port=5433",
		"sslmode=require",
		"TimeZone=Asia/Shanghai",
	} {
		if !strings.Contains(dsn, part) {
			t.Fatalf("dsn missing %q: %s", part, dsn)
		}
	}
}

func TestAutoMigrateEnabled(t *testing.T) {
	orig := os.Getenv("SERVIFY_AUTO_MIGRATE")
	defer os.Setenv("SERVIFY_AUTO_MIGRATE", orig)

	cases := []struct {
		value string
		want  bool
	}{
		{"", true},
		{"true", true},
		{"false", false},
		{"0", false},
		{"1", true},
	}
	for _, tt := range cases {
		if tt.value == "" {
			os.Unsetenv("SERVIFY_AUTO_MIGRATE")
		} else {
			os.Setenv("SERVIFY_AUTO_MIGRATE", tt.value)
		}
		if got := AutoMigrateEnabled(); got != tt.want {
			t.Fatalf("AutoMigrateEnabled(%q)=%v want %v", tt.value, got, tt.want)
		}
	}
}

func TestOpenDatabaseWithRetryFailsAfterAttempts(t *testing.T) {
	start := time.Now()
	_, err := OpenDatabaseWithRetry(config.GetDefaultConfig(), DatabaseOptions{
		DSN: "host=127.0.0.1 user=invalid password=invalid dbname=invalid port=1 sslmode=disable TimeZone=UTC",
	}, DatabaseRetryOptions{
		MaxRetries: 2,
		RetryDelay: time.Nanosecond,
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "connect database after 2 attempts") {
		t.Fatalf("unexpected error: %v", err)
	}
	if errors.Unwrap(err) == nil {
		t.Fatalf("expected wrapped database error: %v", err)
	}
	if time.Since(start) > 5*time.Second {
		t.Fatal("retry test took too long")
	}
}
