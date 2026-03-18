package server

import (
	"net/http"
	"strings"

	"servify/apps/server/internal/config"
	"servify/apps/server/internal/middleware"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func registerBaseMiddleware(r *gin.Engine, cfg *config.Config) {
	if cfg != nil && cfg.Log.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(corsMiddlewareWithConfig(cfg))
	r.Use(middleware.RateLimitMiddlewareFromConfig(cfg))
	if cfg != nil && cfg.Monitoring.Tracing.Enabled {
		r.Use(otelgin.Middleware(cfg.Monitoring.Tracing.ServiceName))
	}
}

func corsMiddlewareWithConfig(cfg *config.Config) gin.HandlerFunc {
	allowedOrigins := "*"
	allowedMethods := "GET, POST, PUT, DELETE, OPTIONS"
	allowedHeaders := "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"
	if cfg != nil && cfg.Security.CORS.Enabled {
		if len(cfg.Security.CORS.AllowedOrigins) > 0 {
			allowedOrigins = strings.Join(cfg.Security.CORS.AllowedOrigins, ", ")
		}
		if len(cfg.Security.CORS.AllowedMethods) > 0 {
			allowedMethods = strings.Join(cfg.Security.CORS.AllowedMethods, ", ")
		}
		if len(cfg.Security.CORS.AllowedHeaders) > 0 {
			allowedHeaders = strings.Join(cfg.Security.CORS.AllowedHeaders, ", ")
		}
	}
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", allowedOrigins)
		c.Header("Access-Control-Allow-Methods", allowedMethods)
		c.Header("Access-Control-Allow-Headers", allowedHeaders)
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
