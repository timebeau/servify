// Package metrics provides Prometheus-based application metrics for Servify.
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

// DefaultRegistry is the process-level Prometheus registry.
// It uses prometheus.DefaultRegisterer so OpenTelemetry and other
// libraries that register against the default registry cooperate.
var DefaultRegistry = NewRegistry()

// Registry wraps a Prometheus registry for metric registration.
type Registry struct {
	registry *prometheus.Registry
}

// NewRegistry creates a new Prometheus registry.
func NewRegistry() *Registry {
	return &Registry{
		registry: prometheus.NewRegistry(),
	}
}

// MustRegister registers a Prometheus Collector. Panics on duplicate registration.
func (r *Registry) MustRegister(c ...prometheus.Collector) {
	r.registry.MustRegister(c...)
}

// Gatherer returns the underlying Gatherer for use with promhttp.
func (r *Registry) Gatherer() prometheus.Gatherer {
	return r.registry
}

// RegisterGoCollector adds the default Go runtime metrics (goroutines, memory, etc).
func (r *Registry) RegisterGoCollector() {
	r.registry.MustRegister(collectors.NewGoCollector())
}

// RegisterProcessCollector adds process-level metrics (CPU, file descriptors, etc).
func (r *Registry) RegisterProcessCollector() {
	r.registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
}
