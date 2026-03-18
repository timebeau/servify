package middleware

import (
	"servify/apps/server/internal/config"
	platformauth "servify/apps/server/internal/platform/auth"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware keeps the legacy middleware entrypoint but delegates to platform/auth.
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return platformauth.AuthMiddleware(platformauth.MiddlewareConfigFromApp(cfg))
}
