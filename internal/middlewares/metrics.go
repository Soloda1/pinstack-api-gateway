package middlewares

import (
	"net/http"
	"pinstack-api-gateway/internal/metrics"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func MetricsMiddleware(metricsProvider metrics.MetricsProvider) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			endpoint := r.URL.Path
			if routeCtx := chi.RouteContext(r.Context()); routeCtx != nil && routeCtx.RoutePattern() != "" {
				endpoint = routeCtx.RoutePattern()
			}

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			metricsProvider.IncHTTPRequestsTotal(r.Method, endpoint, "processing")

			next.ServeHTTP(ww, r)

			duration := time.Since(start)
			status := strconv.Itoa(ww.Status())

			metricsProvider.ObserveHTTPRequestDuration(r.Method, endpoint, duration)
			metricsProvider.IncHTTPResponsesTotal(r.Method, endpoint, status)

			metricsProvider.IncHTTPRequestsTotal(r.Method, endpoint, status)
		})
	}
}
