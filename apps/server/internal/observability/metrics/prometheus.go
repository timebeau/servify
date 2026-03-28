package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// PrometheusHandler returns a Gin handler that serves Prometheus-format metrics.
func PrometheusHandler(reg *Registry) gin.HandlerFunc {
	handler := promhttp.HandlerFor(reg.Gatherer(), promhttp.HandlerOpts{
		EnableOpenMetrics: true,
	})
	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}
