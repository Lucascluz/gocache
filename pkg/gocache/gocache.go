// Package gocache provides a turnkey cache server with HTTP API.
package gocache

import (
	"github.com/Lucascluz/gocache/pkg/cache"
	"github.com/Lucascluz/gocache/pkg/http"
)

// Config holds configuration for cache and HTTP server.
type Config struct {
	CacheConfig cache.Config
	HttpConfig  http.Config
}

// Server wraps a cache with optional HTTP server.
type Server struct {
	cfg     Config
	store   *cache.Cache
	httpSrv *http.Server
}

// New creates a new server with cache and optional HTTP server.
func New(cfg *Config) *Server {
	// init cache
	store := cache.New(&cfg.CacheConfig)

	// init http server if enabled
	var httpSrv *http.Server
	if cfg.HttpConfig.Enabled {
		httpSrv = http.New(&cfg.HttpConfig)
	}

	return &Server{
		cfg:     *cfg,
		store:   store,
		httpSrv: httpSrv,
	}
}

// Cache returns the underlying cache for direct access.
func (s *Server) Cache() *cache.Cache {
	return s.store
}
