package metrics

import (
	"time"
)

type MetricsProvider interface {
	// HTTP Metrics
	IncHTTPRequestsTotal(method, endpoint, status string)
	ObserveHTTPRequestDuration(method, endpoint string, duration time.Duration)
	IncHTTPResponsesTotal(method, endpoint, status string)
	SetActiveHTTPConnections(count int)

	// gRPC Client Metrics
	IncGRPCClientRequestsTotal(service, method, status string)
	ObserveGRPCClientRequestDuration(service, method string, duration time.Duration)
	IncGRPCClientConnectionsTotal(service, status string)

	// Gateway-specific Metrics
	IncProxyRequestsTotal(service, endpoint, status string)
	ObserveProxyRequestDuration(service, endpoint string, duration time.Duration)
	IncAuthenticationTotal(status string)
	IncAuthorizationTotal(endpoint, status string)

	// Rate Limiting Metrics
	IncRateLimitHits(endpoint string)
	IncRateLimitExceeded(endpoint string)

	// Circuit Breaker Metrics
	IncCircuitBreakerStateChanges(service, state string)
	IncCircuitBreakerRequests(service, status string)

	// General System Metrics
	SetActiveUsers(count int)
	IncCacheHits(cache_type string)
	IncCacheMisses(cache_type string)
}
