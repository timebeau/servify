package metrics

import (
	"servify/apps/server/internal/observability/telemetry"

	"github.com/prometheus/client_golang/prometheus"
)

// BusinessMetrics holds counters and histograms for domain operations.
type BusinessMetrics struct {
	conversationsCreated *prometheus.CounterVec
	ticketsCreated       *prometheus.CounterVec
	ticketsResolved      *prometheus.CounterVec
	routingDecisions     *prometheus.CounterVec
	aiRequests           *prometheus.CounterVec
	aiRequestDuration    *prometheus.HistogramVec
	aiLLMTokens          *prometheus.CounterVec
}

// NewBusinessMetrics creates and registers business metric collectors.
func NewBusinessMetrics(reg *Registry) *BusinessMetrics {
	m := &BusinessMetrics{
		conversationsCreated: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: telemetry.MetricConversationsCreated,
			Help: "Total number of conversations created.",
		}, []string{telemetry.LabelTenantID, telemetry.LabelChannel}),
		ticketsCreated: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: telemetry.MetricTicketsCreated,
			Help: "Total number of tickets created.",
		}, []string{telemetry.LabelTenantID, telemetry.LabelPriority}),
		ticketsResolved: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: telemetry.MetricTicketsResolved,
			Help: "Total number of tickets resolved.",
		}, []string{telemetry.LabelTenantID, telemetry.LabelOutcome}),
		routingDecisions: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: telemetry.MetricRoutingDecisions,
			Help: "Total number of routing decisions made.",
		}, []string{telemetry.LabelTenantID, telemetry.LabelStrategy, telemetry.LabelOutcome}),
		aiRequests: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: telemetry.MetricAIRequestsTotal,
			Help: "Total number of AI requests processed.",
		}, []string{telemetry.LabelProvider, telemetry.LabelModel, telemetry.LabelOutcome}),
		aiRequestDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    telemetry.MetricAIRequestDuration,
			Help:    "AI request duration in seconds.",
			Buckets: prometheus.DefBuckets,
		}, []string{telemetry.LabelProvider, telemetry.LabelModel}),
		aiLLMTokens: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: telemetry.MetricAILLMTokenUsage,
			Help: "Total number of LLM tokens consumed.",
		}, []string{telemetry.LabelProvider, telemetry.LabelTokenType}),
	}

	reg.MustRegister(
		m.conversationsCreated,
		m.ticketsCreated,
		m.ticketsResolved,
		m.routingDecisions,
		m.aiRequests,
		m.aiRequestDuration,
		m.aiLLMTokens,
	)
	return m
}

// RecordConversationCreated increments the conversation counter.
func (m *BusinessMetrics) RecordConversationCreated(tenantID, channel string) {
	m.conversationsCreated.WithLabelValues(tenantID, channel).Inc()
}

// RecordTicketCreated increments the ticket creation counter.
func (m *BusinessMetrics) RecordTicketCreated(tenantID, priority string) {
	m.ticketsCreated.WithLabelValues(tenantID, priority).Inc()
}

// RecordTicketResolved increments the ticket resolution counter.
func (m *BusinessMetrics) RecordTicketResolved(tenantID, resolution string) {
	m.ticketsResolved.WithLabelValues(tenantID, resolution).Inc()
}

// RecordRoutingDecision increments routing decisions.
func (m *BusinessMetrics) RecordRoutingDecision(tenantID, strategy, outcome string) {
	m.routingDecisions.WithLabelValues(tenantID, strategy, outcome).Inc()
}

// RecordAIRequest records an AI request with timing.
func (m *BusinessMetrics) RecordAIRequest(provider, model, outcome string, duration float64) {
	m.aiRequests.WithLabelValues(provider, model, outcome).Inc()
	m.aiRequestDuration.WithLabelValues(provider, model).Observe(duration)
}

// RecordAILLMTokens records token consumption.
func (m *BusinessMetrics) RecordAILLMTokens(provider, tokenType string, count int) {
	m.aiLLMTokens.WithLabelValues(provider, tokenType).Add(float64(count))
}
