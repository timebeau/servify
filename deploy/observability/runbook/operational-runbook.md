# Servify Operational Runbook

## Observability Architecture

Servify exposes metrics at `/metrics` in Prometheus format. The stack:

- **Prometheus** scrapes `/metrics` (default every 15s)
- **Jaeger** receives traces via OTel Collector (OTLP gRPC on :4317)
- **Grafana** queries Prometheus for dashboards and alerts

### Metric Naming

All metrics follow Prometheus conventions: `subsystem_name_units`. Key prefixes:

| Prefix | Domain |
|--------|--------|
| `http_` | HTTP request metrics |
| `conversations_` | Conversation events |
| `tickets_` | Ticket lifecycle |
| `routing_` | Routing decisions |
| `ai_` | AI/LLM interactions |
| `eventbus_` | Event bus processing |
| `worker_` | Background job processing |
| `errors_` | Classified errors |

## Alert Runbooks

### HighHTTP5xxRate

**Severity**: Critical | **Threshold**: >5% 5xx rate for 5 minutes

**Investigation**:
1. Check Grafana "Infrastructure" dashboard for which endpoints are failing
2. Check `errors_total{severity="system"}` for classified errors
3. Review application logs filtered by `request_id` from the dashboard
4. Check if a recent deployment coincides with the error spike

**Common causes**:
- Database connection pool exhaustion
- External dependency (LLM, WeKnora) outage
- Configuration error after deployment

**Resolution**:
- Database: check `go_sql_open_connections` metric, increase pool if needed
- External deps: check circuit breaker state, enable fallback mode
- Config: verify environment variables and config.yml

### HighSystemErrorRate

**Severity**: Critical | **Threshold**: system errors > 0.01/s for 5 minutes

**Investigation**:
1. Filter `errors_total` by `error_module` label to find the source
2. Check logs for stack traces matching the module
3. Look for patterns: nil pointer, index out of range, type assertion

**Common causes**:
- Code bug triggered by new request pattern
- Missing validation on new fields
- Concurrency issue under load

### EventBusHandlerFailures

**Severity**: Warning | **Threshold**: any handler failures for 5 minutes

**Investigation**:
1. Check `eventbus_failed_total` by `event_type`
2. Check dead letter entries via the in-memory recorder
3. Review handler code for the failing event type

**Resolution**:
- Transient errors: events will be retried or dead-lettered
- Persistent errors: fix handler code and redeploy
- Use replay interface to reprocess dead-lettered events

### AIProviderDegraded

**Severity**: Critical | **Threshold**: >20% failure rate for 5 minutes

**Investigation**:
1. Check `ai_requests_total` by `provider` and `outcome`
2. Check `ai_request_duration_seconds` for latency spikes
3. Verify API key validity and rate limit status

**Resolution**:
- Rate limited: reduce request rate or increase tier
- Auth failed: rotate API keys
- Timeout: increase timeout config or enable circuit breaker
- Enable fallback to alternative provider

### AIHighLatency

**Severity**: Warning | **Threshold**: P95 latency > 10s for 10 minutes

**Investigation**:
1. Check `ai_request_duration_seconds` by `provider`
2. Compare latency increase with token volume spikes in `ai_llm_tokens_total`
3. Check whether the affected provider is also approaching failure alerts

**Resolution**:
- Reduce concurrent AI load or expensive prompt paths
- Switch to fallback mode when latency makes the product flow unusable
- Validate provider-side latency before increasing timeouts

### WorkerJobFailures

**Severity**: Warning | **Threshold**: any failures for 10 minutes

**Investigation**:
1. Check `worker_jobs_total` by `worker_name` and `outcome`
2. Check `worker_job_duration_seconds` for slow jobs
3. Review worker-specific logs

## Common Operations

### Check Dead Letter Queue

The in-memory dead letter recorder stores up to 1000 entries. Access via internal APIs or application logs filtered by dead letter events.

### Adjust Rate Limits

Edit `config.yml`:
```yaml
security:
  rate_limiting:
    enabled: true
    requests_per_minute: 300
    burst: 50
```

### Enable/Disable Tracing

Edit `config.yml`:
```yaml
monitoring:
  tracing:
    enabled: true
    endpoint: "http://otel-collector:4317"
    sample_ratio: 0.1
```

### Circuit Breaker Management

Reset circuit breaker via API:
```
POST /api/v1/ai/circuit-breaker/reset
```

## Metric Reference

### HTTP Metrics

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `http_requests_total` | Counter | method, path, status_code | Total HTTP requests |
| `http_request_duration_seconds` | Histogram | method, path | Request latency |
| `http_response_size_bytes` | Histogram | method, path | Response size |

### Business Metrics

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `conversations_created_total` | Counter | tenant_id, channel | New conversations |
| `tickets_created_total` | Counter | tenant_id, priority | New tickets |
| `tickets_resolved_total` | Counter | tenant_id, outcome | Resolved tickets |
| `routing_decisions_total` | Counter | tenant_id, strategy, outcome | Routing outcomes |
| `ai_requests_total` | Counter | provider, model, outcome | AI requests |
| `ai_request_duration_seconds` | Histogram | provider, model | AI latency |
| `ai_llm_tokens_total` | Counter | provider, token_type | Token consumption |

### Infrastructure Metrics

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `eventbus_published_total` | Counter | event_type, outcome | Events published |
| `eventbus_handled_total` | Counter | event_type | Events handled |
| `eventbus_failed_total` | Counter | event_type | Handler failures |
| `eventbus_handle_duration_seconds` | Histogram | event_type | Handler duration |
| `eventbus_dead_letter_total` | Counter | event_type | Dead-lettered events |
| `worker_jobs_total` | Counter | worker_name, outcome | Worker jobs |
| `worker_job_duration_seconds` | Histogram | worker_name | Job duration |
| `worker_active_jobs` | Gauge | worker_name | Active jobs |
| `errors_total` | Counter | severity, error_category, error_module | Classified errors |
