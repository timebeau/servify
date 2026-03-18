package server

import (
	"servify/apps/server/internal/handlers"

	"github.com/gin-gonic/gin"
)

func registerHealthRoutes(r routeRegistrar, deps Dependencies) {
	healthHandler := handlers.NewEnhancedHealthHandler(deps.Config, deps.AIService)
	r.GET("/health", healthHandler.Health)
	r.GET("/ready", healthHandler.Ready)

	if deps.Config != nil && deps.Config.Monitoring.Enabled {
		r.GET(deps.Config.Monitoring.MetricsPath, handlers.NewMetricsHandler(
			deps.RealtimeGateway,
			deps.RTCGateway,
			deps.AIService,
			deps.MessageRouter,
			deps.DB,
		).GetMetrics)
	}
}

type routeRegistrar interface {
	GET(string, ...gin.HandlerFunc) gin.IRoutes
}
