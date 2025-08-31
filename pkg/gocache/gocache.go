package gocache

import (
	"github.com/Lucascluz/gocache/pkg/cache"
	"github.com/Lucascluz/gocache/pkg/http"
)

type Config struct {
	CacheConfig cache.Config
	HttpConfig  http.Config
}

type Server struct {
	cfg     Config
	store   *cache.Cache
	httpSrv *http.Server
}

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
