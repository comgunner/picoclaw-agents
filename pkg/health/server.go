package health

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

type Server struct {
	server    *http.Server
	mux       *http.ServeMux
	mu        sync.RWMutex
	ready     bool
	checks    map[string]Check
	startTime time.Time
}

type Check struct {
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	Message   string    `json:"message,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

type StatusResponse struct {
	Status string           `json:"status"`
	Uptime string           `json:"uptime"`
	Pid    int              `json:"pid,omitempty"`
	Checks map[string]Check `json:"checks,omitempty"`
}

func NewServer(host string, port int) *Server {
	mux := http.NewServeMux()
	s := &Server{
		mux:       mux,
		ready:     false,
		checks:    make(map[string]Check),
		startTime: time.Now(),
	}

	mux.HandleFunc("/health", s.healthHandler)
	mux.HandleFunc("/ready", s.readyHandler)

	addr := fmt.Sprintf("%s:%d", host, port)
	s.server = &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	return s
}

// Handle registers the handler for the given pattern on the health server's mux.
// This allows external components (e.g., channel WebSocket handlers) to be mounted
// on the same HTTP port as the health/ready endpoints.
func (s *Server) Handle(pattern string, handler http.Handler) {
	s.mux.Handle(pattern, handler)
}

// HandleFunc registers the handler function for the given pattern on the health server's mux.
func (s *Server) HandleFunc(pattern string, handler http.HandlerFunc) {
	s.mux.HandleFunc(pattern, handler)
}

func (s *Server) Start() error {
	s.mu.Lock()
	s.ready = true
	s.mu.Unlock()
	return s.server.ListenAndServe()
}

func (s *Server) StartContext(ctx context.Context) error {
	s.mu.Lock()
	s.ready = true
	s.mu.Unlock()

	errCh := make(chan error, 1)
	go func() {
		errCh <- s.server.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return s.server.Shutdown(context.Background())
	}
}

func (s *Server) Stop(ctx context.Context) error {
	s.mu.Lock()
	s.ready = false
	s.mu.Unlock()
	return s.server.Shutdown(ctx)
}

func (s *Server) SetReady(ready bool) {
	s.mu.Lock()
	s.ready = ready
	s.mu.Unlock()
}

func (s *Server) RegisterCheck(name string, checkFn func() (bool, string)) {
	s.mu.Lock()
	defer s.mu.Unlock()

	status, msg := checkFn()
	s.checks[name] = Check{
		Name:      name,
		Status:    statusString(status),
		Message:   msg,
		Timestamp: time.Now(),
	}
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	uptime := time.Since(s.startTime)
	resp := StatusResponse{
		Status: "ok",
		Uptime: uptime.String(),
	}

	// BUG-06 FIX: Check JSON encode error
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		// Log encode error but can't change HTTP status (already sent)
		// This helps debugging serialization issues in production
		fmt.Fprintf(os.Stderr, "health: encode response failed: %v\n", err)
	}
}

func (s *Server) readyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	s.mu.RLock()
	ready := s.ready
	checks := make(map[string]Check)
	for k, v := range s.checks {
		checks[k] = v
	}
	s.mu.RUnlock()

	if !ready {
		w.WriteHeader(http.StatusServiceUnavailable)
		// BUG-06 FIX: Check JSON encode error
		if err := json.NewEncoder(w).Encode(StatusResponse{
			Status: "not ready",
			Checks: checks,
		}); err != nil {
			fmt.Fprintf(os.Stderr, "health: encode not-ready response failed: %v\n", err)
		}
		return
	}

	for _, check := range checks {
		if check.Status == "fail" {
			w.WriteHeader(http.StatusServiceUnavailable)
			// BUG-06 FIX: Check JSON encode error
			if err := json.NewEncoder(w).Encode(StatusResponse{
				Status: "not ready",
				Checks: checks,
			}); err != nil {
				fmt.Fprintf(os.Stderr, "health: encode fail response failed: %v\n", err)
			}
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	uptime := time.Since(s.startTime)
	// BUG-06 FIX: Check JSON encode error
	if err := json.NewEncoder(w).Encode(StatusResponse{
		Status: "ready",
		Uptime: uptime.String(),
		Checks: checks,
	}); err != nil {
		fmt.Fprintf(os.Stderr, "health: encode ready response failed: %v\n", err)
	}
}

func statusString(ok bool) string {
	if ok {
		return "ok"
	}
	return "fail"
}
