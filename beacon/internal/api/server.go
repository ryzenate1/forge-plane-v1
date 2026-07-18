package api

import (
	"gamepanel/beacon/internal/auth"
	"gamepanel/beacon/internal/health"
	"gamepanel/beacon/internal/metrics"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// HealthHandler handles health check requests
func HealthHandler(checker health.HealthChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := checker.Check(r.Context()); err != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"status":"unhealthy","error":"` + err.Error() + `"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}
}

func NewServer(
	checker health.HealthChecker,
	metricsCollector metrics.MetricsCollector,
	serverHandler *ServerHandler,
) *http.Server {
	r := mux.NewRouter()

	// Register server CRUD routes
	serverHandler.RegisterRoutes(r)

	// Add health check endpoint at root (for k8s probes etc.)
	r.Handle("/health", HealthHandler(checker)).Methods(http.MethodGet)

	// Add metrics endpoint at root
	r.Handle("/metrics", promhttp.Handler())

	// Log registered routes
	if err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, _ := route.GetPathTemplate()
		methods, _ := route.GetMethods()
		log.Printf("registered route: %s %v", path, methods)
		return nil
	}); err != nil {
		log.Printf("route walk error: %v", err)
	}

	// Apply middleware in correct order (outermost first):
	// Recovery → RequestID → Logging → CORS → CSRF → Auth → RateLimit → ErrorHandler → ValidateRequest
	handler := Chain(
		r,
		RecoveryMiddleware,
		RequestIDMiddleware,
		LoggingMiddleware,
		CORSMiddleware,
		CSRFMiddleware,
		auth.AuthMiddleware,
		RateLimitMiddleware,
		ErrorHandler,
		ValidateRequest,
	)

	// Set up versioned routes under /v1/
	versionedHandlers := []VersionedHandler{
		{
			Version: CurrentVersion,
			Handler: handler,
		},
	}
	versionedRouter := VersionedRouter(versionedHandlers...)

	// Root-level mux ensures /health is accessible without /v1/ prefix
	rootMux := http.NewServeMux()
	rootMux.Handle("/health", HealthHandler(checker))
	rootMux.Handle("/", versionedRouter)

	return &http.Server{
		Handler:           rootMux,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}
