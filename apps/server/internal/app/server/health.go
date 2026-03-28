package server

import (
	"servify/apps/server/internal/handlers"
	svcmetrics "servify/apps/server/internal/observability/metrics"

	"github.com/gin-gonic/gin"
)

func registerHealthRoutes(r routeRegistrar, deps Dependencies) {
	healthHandler := handlers.NewEnhancedHealthHandler(deps.Config, deps.AIHandlerService)
	r.GET("/health", healthHandler.Health)
	r.GET("/ready", healthHandler.Ready)

	if deps.Config != nil && deps.Config.Monitoring.Enabled {
		if deps.HTTPMetrics != nil {
			r.GET(deps.Config.Monitoring.MetricsPath, svcmetrics.PrometheusHandler(svcmetrics.DefaultRegistry))
		} else {
			r.GET(deps.Config.Monitoring.MetricsPath, handlers.NewMetricsHandler(
				deps.RealtimeGateway,
				deps.RTCGateway,
				deps.AIHandlerService,
				newGormDBStatsProvider(deps.DB),
			).GetMetrics)
		}
	}
}

type routeRegistrar interface {
	GET(string, ...gin.HandlerFunc) gin.IRoutes
}
