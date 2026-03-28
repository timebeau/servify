package errors

import (
	"servify/apps/server/internal/observability/metrics"
	"servify/apps/server/internal/observability/telemetry"

	"github.com/prometheus/client_golang/prometheus"
)

var errorTotal *prometheus.CounterVec

// RegisterErrorMetrics creates and registers the error counter.
func RegisterErrorMetrics(reg *metrics.Registry) {
	errorTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: telemetry.MetricErrorsTotal,
		Help: "Total number of classified errors.",
	}, []string{telemetry.LabelSeverity, telemetry.LabelErrorCategory, telemetry.LabelErrorModule})

	reg.MustRegister(errorTotal)
}

// RecordError increments the error counter for the given AppError classification.
func RecordError(appErr *AppError) {
	if errorTotal == nil || appErr == nil {
		return
	}
	errorTotal.WithLabelValues(
		string(appErr.Severity),
		string(appErr.Category),
		appErr.Module,
	).Inc()
}
