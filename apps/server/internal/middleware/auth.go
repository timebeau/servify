package middleware

import (
	"servify/apps/server/internal/config"
	platformauth "servify/apps/server/internal/platform/auth"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware keeps the legacy middleware entrypoint but delegates to platform/auth.
func AuthMiddleware(cfg *config.Config, policies ...platformauth.TokenPolicy) gin.HandlerFunc {
	mwCfg := platformauth.MiddlewareConfigFromApp(cfg)
	mwCfg.Policy = platformauth.ComposeTokenPolicies(policies...)
	return platformauth.AuthMiddleware(mwCfg)
}

// EnforceRequestScope keeps the compatibility layer aligned with platform/auth scope rules.
func EnforceRequestScope() gin.HandlerFunc {
	return platformauth.EnforceRequestScope()
}
