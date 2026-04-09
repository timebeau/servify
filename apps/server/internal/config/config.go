package config

import (
	"time"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

type Config struct {
	Server     ServerConfig     `yaml:"server"`
	Database   DatabaseConfig   `yaml:"database"`
	Redis      RedisConfig      `yaml:"redis"`
	WebRTC     WebRTCConfig     `yaml:"webrtc"`
	AI         AIConfig         `yaml:"ai"`
	WeKnora    WeKnoraConfig    `yaml:"weknora"`
	Fallback   FallbackConfig   `yaml:"fallback"`
	JWT        JWTConfig        `yaml:"jwt"`
	Log        LogConfig        `yaml:"log"`
	Monitoring MonitoringConfig `yaml:"monitoring"`
	Security   SecurityConfig   `yaml:"security"`
	Portal     PortalConfig     `yaml:"portal"`
	Upload     UploadConfig     `yaml:"upload"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type DatabaseConfig struct {
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	User            string        `yaml:"user"`
	Password        string        `yaml:"password"`
	Name            string        `yaml:"name"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
}

type RedisConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	Password     string `yaml:"password"`
	DB           int    `yaml:"db"`
	PoolSize     int    `yaml:"pool_size"`
	MinIdleConns int    `yaml:"min_idle_conns"`
}

type WebRTCConfig struct {
	STUNServer string `yaml:"stun_server"`
}

type AIConfig struct {
	OpenAI OpenAIConfig `yaml:"openai"`
}

type OpenAIConfig struct {
	APIKey      string        `yaml:"api_key" json:"api_key,omitempty"`
	BaseURL     string        `yaml:"base_url" json:"base_url,omitempty"`
	Model       string        `yaml:"model" json:"model,omitempty"`
	Temperature float64       `yaml:"temperature" json:"temperature,omitempty"`
	MaxTokens   int           `yaml:"max_tokens" json:"max_tokens,omitempty"`
	Timeout     time.Duration `yaml:"timeout" json:"timeout,omitempty"`
}

type WeKnoraConfig struct {
	Enabled         bool                `yaml:"enabled" json:"enabled,omitempty"`
	BaseURL         string              `yaml:"base_url" json:"base_url,omitempty"`
	APIKey          string              `yaml:"api_key" json:"api_key,omitempty"`
	TenantID        string              `yaml:"tenant_id" json:"tenant_id,omitempty"`
	KnowledgeBaseID string              `yaml:"knowledge_base_id" json:"knowledge_base_id,omitempty"`
	Timeout         time.Duration       `yaml:"timeout" json:"timeout,omitempty"`
	MaxRetries      int                 `yaml:"max_retries" json:"max_retries,omitempty"`
	Search          WeKnoraSearchConfig `yaml:"search" json:"search,omitempty"`
	HealthCheck     WeKnoraHealthConfig `yaml:"health_check" json:"health_check,omitempty"`
}

type WeKnoraSearchConfig struct {
	DefaultLimit   int     `yaml:"default_limit" json:"default_limit,omitempty"`
	ScoreThreshold float64 `yaml:"score_threshold" json:"score_threshold,omitempty"`
	Strategy       string  `yaml:"strategy" json:"strategy,omitempty"`
}

type WeKnoraHealthConfig struct {
	Interval time.Duration `yaml:"interval" json:"interval,omitempty"`
	Timeout  time.Duration `yaml:"timeout" json:"timeout,omitempty"`
}

type FallbackConfig struct {
	Enabled         bool                 `yaml:"enabled"`
	LegacyKBEnabled bool                 `yaml:"legacy_kb_enabled"`
	CircuitBreaker  CircuitBreakerConfig `yaml:"circuit_breaker"`
}

type CircuitBreakerConfig struct {
	Enabled         bool          `yaml:"enabled"`
	MaxFailures     int           `yaml:"max_failures"`
	ResetTimeout    time.Duration `yaml:"reset_timeout"`
	HalfOpenMaxReqs int           `yaml:"half_open_max_requests"`
}

type JWTConfig struct {
	Secret           string        `yaml:"secret"`
	ExpiresIn        time.Duration `yaml:"expires_in"`
	RefreshExpiresIn time.Duration `yaml:"refresh_expires_in"`
}

type LogConfig struct {
	Level      string `yaml:"level"`
	Format     string `yaml:"format"` // json, text
	Output     string `yaml:"output"` // stdout, file, both
	FilePath   string `yaml:"file_path"`
	MaxSize    int    `yaml:"max_size"`    // MB
	MaxAge     int    `yaml:"max_age"`     // days
	MaxBackups int    `yaml:"max_backups"` // number of backup files
	Compress   bool   `yaml:"compress"`    // compress backup files
}

type MonitoringConfig struct {
	Enabled      bool                     `yaml:"enabled"`
	MetricsPath  string                   `yaml:"metrics_path"`
	Performance  PerformanceMonitorConfig `yaml:"performance"`
	HealthChecks HealthChecksConfig       `yaml:"health_checks"`
	Tracing      TracingConfig            `yaml:"tracing"`
}

type PerformanceMonitorConfig struct {
	SlowQueryThreshold   time.Duration `yaml:"slow_query_threshold"`
	EnableRequestLogging bool          `yaml:"enable_request_logging"`
}

type HealthChecksConfig struct {
	Database bool `yaml:"database"`
	Redis    bool `yaml:"redis"`
	WeKnora  bool `yaml:"weknora"`
	OpenAI   bool `yaml:"openai"`
}

// TracingConfig OpenTelemetry 追踪配置
type TracingConfig struct {
	Enabled     bool    `yaml:"enabled"`
	Endpoint    string  `yaml:"endpoint"`     // OTLP gRPC 端点，例如 http://otel-collector:4317 或 0.0.0.0:4317
	Insecure    bool    `yaml:"insecure"`     // 是否使用明文（本地/开发）
	SampleRatio float64 `yaml:"sample_ratio"` // 采样率 0.0~1.0
	ServiceName string  `yaml:"service_name"` // 自定义服务名，缺省使用 "servify"
}

type SecurityConfig struct {
	CORS            CORSConfig              `yaml:"cors"`
	RateLimiting    RateLimitingConfig      `yaml:"rate_limiting"`
	RBAC            RBACConfig              `yaml:"rbac"`
	Audit           AuditConfig             `yaml:"audit"`
	TokenRevocation TokenRevocationConfig   `yaml:"token_revocation"`
	SessionRisk     SessionRiskPolicyConfig `yaml:"session_risk"`
}

type CORSConfig struct {
	Enabled        bool     `yaml:"enabled"`
	AllowedOrigins []string `yaml:"allowed_origins"`
	AllowedMethods []string `yaml:"allowed_methods"`
	AllowedHeaders []string `yaml:"allowed_headers"`
}

type RBACConfig struct {
	Enabled bool                `yaml:"enabled"`
	Roles   map[string][]string `yaml:"roles"`
}

type AuditConfig struct {
	Enabled          bool          `yaml:"enabled"`
	Retention        time.Duration `yaml:"retention"`
	CleanupInterval  time.Duration `yaml:"cleanup_interval"`
	CleanupBatchSize int           `yaml:"cleanup_batch_size"`
}

type TokenRevocationConfig struct {
	Enabled          bool          `yaml:"enabled"`
	CleanupInterval  time.Duration `yaml:"cleanup_interval"`
	CleanupBatchSize int           `yaml:"cleanup_batch_size"`
}

type SessionRiskPolicyConfig struct {
	HotRefreshWindowMinutes    int `yaml:"hot_refresh_window_minutes" json:"hot_refresh_window_minutes,omitempty"`
	RecentRefreshWindowMinutes int `yaml:"recent_refresh_window_minutes" json:"recent_refresh_window_minutes,omitempty"`
	TodayRefreshWindowHours    int `yaml:"today_refresh_window_hours" json:"today_refresh_window_hours,omitempty"`
	RapidChangeWindowHours     int `yaml:"rapid_change_window_hours" json:"rapid_change_window_hours,omitempty"`
	StaleActivityWindowDays    int `yaml:"stale_activity_window_days" json:"stale_activity_window_days,omitempty"`
	MultiPublicIPThreshold     int `yaml:"multi_public_ip_threshold" json:"multi_public_ip_threshold,omitempty"`
	ManySessionsThreshold      int `yaml:"many_sessions_threshold" json:"many_sessions_threshold,omitempty"`
	HotRefreshFamilyThreshold  int `yaml:"hot_refresh_family_threshold" json:"hot_refresh_family_threshold,omitempty"`
	MediumRiskScore            int `yaml:"medium_risk_score" json:"medium_risk_score,omitempty"`
	HighRiskScore              int `yaml:"high_risk_score" json:"high_risk_score,omitempty"`
}

type RateLimitingConfig struct {
	Enabled           bool                  `yaml:"enabled"`
	RequestsPerMinute int                   `yaml:"requests_per_minute"`
	Burst             int                   `yaml:"burst"`
	Paths             []PathRateLimitConfig `yaml:"paths"`
	// Optional: use specific header value as rate-limit key (e.g., X-Forwarded-For, X-API-Key)
	KeyHeader string `yaml:"key_header"`
	// Optional: bypass limit for these IPs (matches client IP)
	WhitelistIPs []string `yaml:"whitelist_ips"`
	// Optional: bypass limit for these header key values (when KeyHeader set)
	WhitelistKeys []string `yaml:"whitelist_keys"`
}

// PathRateLimitConfig allows overriding rate limits for specific path prefixes.
// The first matching Prefix will be used.
type PathRateLimitConfig struct {
	Enabled           bool   `yaml:"enabled"`
	Prefix            string `yaml:"prefix"`
	RequestsPerMinute int    `yaml:"requests_per_minute"`
	Burst             int    `yaml:"burst"`
}

// PortalConfig controls public portal branding and i18n defaults for static pages.
type PortalConfig struct {
	BrandName      string   `yaml:"brand_name" json:"brand_name,omitempty"`
	LogoURL        string   `yaml:"logo_url" json:"logo_url,omitempty"`
	PrimaryColor   string   `yaml:"primary_color" json:"primary_color,omitempty"`
	SecondaryColor string   `yaml:"secondary_color" json:"secondary_color,omitempty"`
	DefaultLocale  string   `yaml:"default_locale" json:"default_locale,omitempty"` // e.g. zh-CN, en-US
	Locales        []string `yaml:"locales" json:"locales,omitempty"`               // allowed locales
	SupportEmail   string   `yaml:"support_email" json:"support_email,omitempty"`
}
type UploadConfig struct {
	Enabled      bool     `yaml:"enabled"`
	MaxFileSize  string   `yaml:"max_file_size"`
	AllowedTypes []string `yaml:"allowed_types"`
	StoragePath  string   `yaml:"storage_path"`
	AutoProcess  bool     `yaml:"auto_process"`
	AutoIndex    bool     `yaml:"auto_index"`
}

func Load() *Config {
	// Start with default config to ensure all fields have valid defaults
	config := GetDefaultConfig()
	// Viper unmarshalling uses mapstructure tags by default; explicitly decode via our `yaml` tags
	// to keep config files consistent (e.g. `stun_server`, `max_open_conns`, etc.).
	if err := viper.Unmarshal(config, func(dc *mapstructure.DecoderConfig) {
		dc.TagName = "yaml"
	}); err != nil {
		panic(err)
	}
	return config
}

// GetDefaultConfig 返回默认配置
func GetDefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Host: "0.0.0.0",
			Port: 8080,
		},
		Database: DatabaseConfig{
			Host:            "localhost",
			Port:            5432,
			User:            "postgres",
			Password:        "password",
			Name:            "servify",
			MaxOpenConns:    100,
			MaxIdleConns:    10,
			ConnMaxLifetime: 3600 * time.Second,
		},
		Redis: RedisConfig{
			Host:         "localhost",
			Port:         6379,
			Password:     "",
			DB:           0,
			PoolSize:     10,
			MinIdleConns: 5,
		},
		WebRTC: WebRTCConfig{
			STUNServer: "stun:stun.l.google.com:19302",
		},
		AI: AIConfig{
			OpenAI: OpenAIConfig{
				BaseURL:     "https://api.openai.com/v1",
				Model:       "gpt-3.5-turbo",
				Temperature: 0.7,
				MaxTokens:   1000,
				Timeout:     30 * time.Second,
			},
		},
		WeKnora: WeKnoraConfig{
			Enabled:         false,
			BaseURL:         "http://localhost:9000",
			APIKey:          "default-api-key",
			TenantID:        "default-tenant",
			KnowledgeBaseID: "default-kb",
			Timeout:         30 * time.Second,
			MaxRetries:      3,
			Search: WeKnoraSearchConfig{
				DefaultLimit:   5,
				ScoreThreshold: 0.7,
				Strategy:       "hybrid",
			},
			HealthCheck: WeKnoraHealthConfig{
				Interval: 30 * time.Second,
				Timeout:  10 * time.Second,
			},
		},
		Fallback: FallbackConfig{
			Enabled:         true,
			LegacyKBEnabled: true,
			CircuitBreaker: CircuitBreakerConfig{
				Enabled:         true,
				MaxFailures:     5,
				ResetTimeout:    60 * time.Second,
				HalfOpenMaxReqs: 3,
			},
		},
		JWT: JWTConfig{
			Secret:           "default-secret-key",
			ExpiresIn:        24 * time.Hour,
			RefreshExpiresIn: 7 * 24 * time.Hour,
		},
		Log: LogConfig{
			Level:      "info",
			Format:     "json",
			Output:     "both",
			FilePath:   "./logs/servify.log",
			MaxSize:    100,
			MaxAge:     7,
			MaxBackups: 3,
			Compress:   true,
		},
		Monitoring: MonitoringConfig{
			Enabled:     true,
			MetricsPath: "/metrics",
			Performance: PerformanceMonitorConfig{
				SlowQueryThreshold:   1 * time.Second,
				EnableRequestLogging: true,
			},
			HealthChecks: HealthChecksConfig{
				Database: true,
				Redis:    true,
				WeKnora:  true,
				OpenAI:   false,
			},
			Tracing: TracingConfig{
				Enabled:     false,
				Endpoint:    "http://localhost:4317",
				Insecure:    true,
				SampleRatio: 0.1,
				ServiceName: "servify",
			},
		},
		Security: SecurityConfig{
			CORS: CORSConfig{
				Enabled:        true,
				AllowedOrigins: []string{"*"},
				AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
				AllowedHeaders: []string{"*"},
			},
			RateLimiting: RateLimitingConfig{
				Enabled:           false,
				RequestsPerMinute: 300,
				Burst:             50,
			},
			Audit: AuditConfig{
				Enabled:          true,
				Retention:        180 * 24 * time.Hour,
				CleanupInterval:  24 * time.Hour,
				CleanupBatchSize: 500,
			},
			TokenRevocation: TokenRevocationConfig{
				Enabled:          true,
				CleanupInterval:  24 * time.Hour,
				CleanupBatchSize: 500,
			},
			SessionRisk: SessionRiskPolicyConfig{
				HotRefreshWindowMinutes:    15,
				RecentRefreshWindowMinutes: 60,
				TodayRefreshWindowHours:    24,
				RapidChangeWindowHours:     24,
				StaleActivityWindowDays:    30,
				MultiPublicIPThreshold:     2,
				ManySessionsThreshold:      3,
				HotRefreshFamilyThreshold:  2,
				MediumRiskScore:            2,
				HighRiskScore:              4,
			},
		},
		Portal: PortalConfig{
			BrandName:      "Servify",
			LogoURL:        "",
			PrimaryColor:   "#4299e1",
			SecondaryColor: "#764ba2",
			DefaultLocale:  "zh-CN",
			Locales:        []string{"zh-CN", "en-US"},
			SupportEmail:   "",
		},
		Upload: UploadConfig{
			Enabled:      true,
			MaxFileSize:  "10MB",
			AllowedTypes: []string{".pdf", ".docx", ".txt", ".md", ".png", ".jpg", ".jpeg"},
			StoragePath:  "./.runtime/uploads",
			AutoProcess:  true,
			AutoIndex:    true,
		},
	}
}
