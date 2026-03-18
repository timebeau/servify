package bootstrap

import (
	"os"
	"strings"
	"testing"

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
