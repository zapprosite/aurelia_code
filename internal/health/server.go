package health

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/kocar/aurelia/internal/observability"
)

// HealthStatus represents the current health status of the system
type HealthStatus struct {
	Status    string                 `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Uptime    string                 `json:"uptime"`
	Checks    map[string]CheckResult `json:"checks"`
}

// CheckResult represents the result of a health check
type CheckResult struct {
	Status  string `json:"status"` // "ok", "warning", or "error"
	Message string `json:"message,omitempty"`
}

// Server provides HTTP health endpoints
type Server struct {
	port      int
	logger    *slog.Logger
	startTime time.Time
	mu        sync.RWMutex
	checks    map[string]func() CheckResult
	routes    map[string]http.Handler
	server    *http.Server
}

// NewServer creates a new health server
func NewServer(port int) *Server {
	return &Server{
		port:      port,
		logger:    observability.Logger("health"),
		startTime: time.Now(),
		checks:    make(map[string]func() CheckResult),
		routes:    make(map[string]http.Handler),
	}
}

// RegisterCheck registers a health check function
func (s *Server) RegisterCheck(name string, fn func() CheckResult) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.checks[name] = fn
}

// RegisterRoute registers an additional HTTP route on the internal server.
func (s *Server) RegisterRoute(path string, handler http.Handler) {
	if path == "" || handler == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.routes[path] = handler
}

// Start starts the health server
func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/ready", s.handleReady)

	s.mu.RLock()
	for path, handler := range s.routes {
		mux.Handle(path, handler)
	}
	s.mu.RUnlock()

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: mux,
	}

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("Health server error", slog.Any("err", err))
		}
	}()

	s.logger.Info("Health server started", slog.Int("port", s.port))
	return nil
}

// Stop stops the health server
func (s *Server) Stop() error {
	if s.server == nil {
		return nil
	}
	s.logger.Info("Stopping health server")
	return s.server.Close()
}

// handleHealth handles /health endpoint
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	checks := make(map[string]func() CheckResult)
	for name, fn := range s.checks {
		checks[name] = fn
	}
	s.mu.RUnlock()

	status := "ok"
	results := make(map[string]CheckResult)

	for name, fn := range checks {
		result := fn()
		results[name] = result
		if result.Status == "error" {
			status = "degraded"
		}
	}

	uptime := time.Since(s.startTime).String()

	response := HealthStatus{
		Status:    status,
		Timestamp: time.Now(),
		Uptime:    uptime,
		Checks:    results,
	}

	w.Header().Set("Content-Type", "application/json")
	if status == "ok" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	json.NewEncoder(w).Encode(response)
}

// handleReady handles /ready endpoint (simplified readiness check)
func (s *Server) handleReady(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status": "ready",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
