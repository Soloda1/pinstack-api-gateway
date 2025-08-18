package middlewares

import (
	"net/http"
	"pinstack-api-gateway/internal/metrics"

	"github.com/go-chi/chi/v5"
)

// AuthMetricsMiddleware returns a middleware that collects authorization metrics
func AuthMetricsMiddleware(metricsProvider metrics.MetricsProvider) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			endpoint := r.URL.Path
			if routeCtx := chi.RouteContext(r.Context()); routeCtx != nil && routeCtx.RoutePattern() != "" {
				endpoint = routeCtx.RoutePattern()
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				metricsProvider.IncAuthorizationTotal(endpoint, "unauthorized")
				next.ServeHTTP(w, r)
				return
			}

			metricsProvider.IncAuthorizationTotal(endpoint, "authorized")
			next.ServeHTTP(w, r)
		})
	}
}
