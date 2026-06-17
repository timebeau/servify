package server

import (
	"strings"
	"time"

	"servify/apps/server/internal/config"
	"servify/apps/server/internal/handlers"
	"servify/apps/server/internal/middleware"
	platformauth "servify/apps/server/internal/platform/auth"
	"servify/apps/server/internal/platform/configscope"
	"servify/apps/server/internal/platform/storage"
	"servify/apps/server/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func authPolicies(db *gorm.DB) []platformauth.TokenPolicy {
	return []platformauth.TokenPolicy{
		platformauth.NewRevokedTokenPolicy(db),
		platformauth.NewUserStateTokenPolicy(db),
	}
}

func registerAuthRoutes(r *gin.Engine, deps Dependencies) {
	auth := r.Group("/api/v1/auth")
	sessionRiskResolver := configscope.NewResolver(
		deps.Config,
		configscope.WithTenantSessionRiskProvider(configscope.NewGormTenantConfigProvider(deps.DB)),
		configscope.WithWorkspaceSessionRiskProvider(configscope.NewGormWorkspaceConfigProvider(deps.DB)),
	)
	authHandler := handlers.NewAuthHandler(services.NewAuthService(deps.DB, deps.Config)).WithSessionRiskResolver(sessionRiskResolver)
	if provider := sessionIPIntelligenceFromConfig(deps.Config); provider != nil {
		authHandler = authHandler.WithSessionIPIntelligence(provider)
	}

	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
	auth.POST("/refresh", authHandler.RefreshToken)

	authMe := auth.Group("")
	authMe.Use(middleware.AuthMiddleware(deps.Config, authPolicies(deps.DB)...))
	authMe.GET("/me", authHandler.GetCurrentUser)
	authMe.GET("/sessions", authHandler.ListSessions)
	authMe.POST("/sessions/logout-current", authHandler.LogoutCurrentSession)
	authMe.POST("/sessions/logout-others", authHandler.LogoutOtherSessions)

	localStorage := storage.NewLocalProvider("./uploads", "/uploads")
	uploadHandler := handlers.NewFileUploadHandler(localStorage, 32<<20)
	r.POST("/api/v1/upload", middleware.AuthMiddleware(deps.Config, authPolicies(deps.DB)...), uploadHandler.Upload)
	r.Static("/uploads", "./uploads")
}

func sessionIPIntelligenceFromConfig(cfg *config.Config) *handlers.HTTPSessionIPIntelligence {
	if cfg == nil {
		return nil
	}
	providerCfg := cfg.Security.SessionIPIntelligence
	if !providerCfg.Enabled || strings.TrimSpace(providerCfg.BaseURL) == "" {
		return nil
	}
	timeout := time.Duration(providerCfg.TimeoutMs) * time.Millisecond
	return handlers.NewHTTPSessionIPIntelligence(providerCfg.BaseURL, providerCfg.APIKey, providerCfg.AuthHeader, timeout)
}
