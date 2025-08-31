package http

import (
	"context"
	stdhttp "net/http"
	"strconv"
)

type Config struct {
	Enabled bool
	Port    int
}

type Server struct {
	cfg Config
	srv *stdhttp.Server
	mux *stdhttp.ServeMux
}

func New(cfg *Config) *Server {
	// Default configuration if none provided
	if cfg == nil {
		cfg = &Config{Enabled: true, Port: 8080}
	}

	mux := stdhttp.NewServeMux()
	srv := &stdhttp.Server{
		Addr:    ":" + strconv.Itoa(cfg.Port),
		Handler: mux,
	}

	s := &Server{cfg: *cfg, srv: srv, mux: mux}

	// Basic health endpoint
	s.mux.HandleFunc("/health", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		w.WriteHeader(stdhttp.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	return s
}

// Handle registers a handler for a pattern.
func (s *Server) Handle(pattern string, handler stdhttp.Handler) { s.mux.Handle(pattern, handler) }

// HandleFunc registers a handler function for a pattern.
func (s *Server) HandleFunc(pattern string, handler func(stdhttp.ResponseWriter, *stdhttp.Request)) {
	s.mux.HandleFunc(pattern, handler)
}

// ListenAndServe starts the server and blocks.
func (s *Server) ListenAndServe() error { return s.srv.ListenAndServe() }

// Shutdown gracefully stops the server.
func (s *Server) Shutdown(ctx context.Context) error { return s.srv.Shutdown(ctx) }

// Addr returns the listen address.
func (s *Server) Addr() string { return s.srv.Addr }
