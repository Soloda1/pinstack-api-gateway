package prometheus

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type PrometheusMetrics struct {
	// HTTP Metrics
	httpRequestsTotal     *prometheus.CounterVec
	httpRequestDuration   *prometheus.HistogramVec
	httpResponsesTotal    *prometheus.CounterVec
	activeHTTPConnections prometheus.Gauge

	// gRPC Client Metrics
	grpcClientRequestsTotal    *prometheus.CounterVec
	grpcClientRequestDuration  *prometheus.HistogramVec
	grpcClientConnectionsTotal *prometheus.CounterVec

	// Gateway-specific Metrics
	proxyRequestsTotal   *prometheus.CounterVec
	proxyRequestDuration *prometheus.HistogramVec
	authenticationTotal  *prometheus.CounterVec
	authorizationTotal   *prometheus.CounterVec

	// Rate Limiting Metrics
	rateLimitHits     *prometheus.CounterVec
	rateLimitExceeded *prometheus.CounterVec

	// Circuit Breaker Metrics
	circuitBreakerStateChanges *prometheus.CounterVec
	circuitBreakerRequests     *prometheus.CounterVec

	// General System Metrics
	activeUsers prometheus.Gauge
	cacheHits   *prometheus.CounterVec
	cacheMisses *prometheus.CounterVec
}

func NewPrometheusMetrics() *PrometheusMetrics {
	metrics := &PrometheusMetrics{
		// HTTP Metrics
		httpRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests received by the API Gateway",
			},
			[]string{"method", "endpoint", "status"},
		),
		httpRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Duration of HTTP requests in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "endpoint"},
		),
		httpResponsesTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_responses_total",
				Help: "Total number of HTTP responses sent by the API Gateway",
			},
			[]string{"method", "endpoint", "status"},
		),
		activeHTTPConnections: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "active_http_connections",
				Help: "Number of active HTTP connections",
			},
		),

		// gRPC Client Metrics
		grpcClientRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "grpc_client_requests_total",
				Help: "Total number of gRPC client requests made by the API Gateway",
			},
			[]string{"service", "method", "status"},
		),
		grpcClientRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "grpc_client_request_duration_seconds",
				Help:    "Duration of gRPC client requests in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"service", "method"},
		),
		grpcClientConnectionsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "grpc_client_connections_total",
				Help: "Total number of gRPC client connection attempts",
			},
			[]string{"service", "status"},
		),

		// Gateway-specific Metrics
		proxyRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "proxy_requests_total",
				Help: "Total number of proxy requests to backend services",
			},
			[]string{"service", "endpoint", "status"},
		),
		proxyRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "proxy_request_duration_seconds",
				Help:    "Duration of proxy requests to backend services in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"service", "endpoint"},
		),
		authenticationTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "authentication_total",
				Help: "Total number of authentication attempts",
			},
			[]string{"status"},
		),
		authorizationTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "authorization_total",
				Help: "Total number of authorization attempts",
			},
			[]string{"endpoint", "status"},
		),

		// Rate Limiting Metrics
		rateLimitHits: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "rate_limit_hits_total",
				Help: "Total number of rate limit hits",
			},
			[]string{"endpoint"},
		),
		rateLimitExceeded: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "rate_limit_exceeded_total",
				Help: "Total number of rate limit exceeded events",
			},
			[]string{"endpoint"},
		),

		// Circuit Breaker Metrics
		circuitBreakerStateChanges: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "circuit_breaker_state_changes_total",
				Help: "Total number of circuit breaker state changes",
			},
			[]string{"service", "state"},
		),
		circuitBreakerRequests: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "circuit_breaker_requests_total",
				Help: "Total number of requests to circuit breaker",
			},
			[]string{"service", "status"},
		),

		// General System Metrics
		activeUsers: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "active_users",
				Help: "Number of active users connected to the API Gateway",
			},
		),
		cacheHits: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "cache_hits_total",
				Help: "Total number of cache hits",
			},
			[]string{"cache_type"},
		),
		cacheMisses: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "cache_misses_total",
				Help: "Total number of cache misses",
			},
			[]string{"cache_type"},
		),
	}

	// Register all metrics with Prometheus
	prometheus.MustRegister(
		metrics.httpRequestsTotal,
		metrics.httpRequestDuration,
		metrics.httpResponsesTotal,
		metrics.activeHTTPConnections,
		metrics.grpcClientRequestsTotal,
		metrics.grpcClientRequestDuration,
		metrics.grpcClientConnectionsTotal,
		metrics.proxyRequestsTotal,
		metrics.proxyRequestDuration,
		metrics.authenticationTotal,
		metrics.authorizationTotal,
		metrics.rateLimitHits,
		metrics.rateLimitExceeded,
		metrics.circuitBreakerStateChanges,
		metrics.circuitBreakerRequests,
		metrics.activeUsers,
		metrics.cacheHits,
		metrics.cacheMisses,
	)

	return metrics
}

// HTTP Metrics Implementation
func (p *PrometheusMetrics) IncHTTPRequestsTotal(method, endpoint, status string) {
	p.httpRequestsTotal.WithLabelValues(method, endpoint, status).Inc()
}

func (p *PrometheusMetrics) ObserveHTTPRequestDuration(method, endpoint string, duration time.Duration) {
	p.httpRequestDuration.WithLabelValues(method, endpoint).Observe(duration.Seconds())
}

func (p *PrometheusMetrics) IncHTTPResponsesTotal(method, endpoint, status string) {
	p.httpResponsesTotal.WithLabelValues(method, endpoint, status).Inc()
}

func (p *PrometheusMetrics) SetActiveHTTPConnections(count int) {
	p.activeHTTPConnections.Set(float64(count))
}

// gRPC Client Metrics Implementation
func (p *PrometheusMetrics) IncGRPCClientRequestsTotal(service, method, status string) {
	p.grpcClientRequestsTotal.WithLabelValues(service, method, status).Inc()
}

func (p *PrometheusMetrics) ObserveGRPCClientRequestDuration(service, method string, duration time.Duration) {
	p.grpcClientRequestDuration.WithLabelValues(service, method).Observe(duration.Seconds())
}

func (p *PrometheusMetrics) IncGRPCClientConnectionsTotal(service, status string) {
	p.grpcClientConnectionsTotal.WithLabelValues(service, status).Inc()
}

// Gateway-specific Metrics Implementation
func (p *PrometheusMetrics) IncProxyRequestsTotal(service, endpoint, status string) {
	p.proxyRequestsTotal.WithLabelValues(service, endpoint, status).Inc()
}

func (p *PrometheusMetrics) ObserveProxyRequestDuration(service, endpoint string, duration time.Duration) {
	p.proxyRequestDuration.WithLabelValues(service, endpoint).Observe(duration.Seconds())
}

func (p *PrometheusMetrics) IncAuthenticationTotal(status string) {
	p.authenticationTotal.WithLabelValues(status).Inc()
}

func (p *PrometheusMetrics) IncAuthorizationTotal(endpoint, status string) {
	p.authorizationTotal.WithLabelValues(endpoint, status).Inc()
}

// Rate Limiting Metrics Implementation
func (p *PrometheusMetrics) IncRateLimitHits(endpoint string) {
	p.rateLimitHits.WithLabelValues(endpoint).Inc()
}

func (p *PrometheusMetrics) IncRateLimitExceeded(endpoint string) {
	p.rateLimitExceeded.WithLabelValues(endpoint).Inc()
}

// Circuit Breaker Metrics Implementation
func (p *PrometheusMetrics) IncCircuitBreakerStateChanges(service, state string) {
	p.circuitBreakerStateChanges.WithLabelValues(service, state).Inc()
}

func (p *PrometheusMetrics) IncCircuitBreakerRequests(service, status string) {
	p.circuitBreakerRequests.WithLabelValues(service, status).Inc()
}

// General System Metrics Implementation
func (p *PrometheusMetrics) SetActiveUsers(count int) {
	p.activeUsers.Set(float64(count))
}

func (p *PrometheusMetrics) IncCacheHits(cacheType string) {
	p.cacheHits.WithLabelValues(cacheType).Inc()
}

func (p *PrometheusMetrics) IncCacheMisses(cacheType string) {
	p.cacheMisses.WithLabelValues(cacheType).Inc()
}
