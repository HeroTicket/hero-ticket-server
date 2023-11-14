package http

import (
	"context"
	"net/http"
	"time"
)

var (
	ErrServerClosed = http.ErrServerClosed
)

type Config struct {
	Version      string
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func DefaultConfig() *Config {
	return &Config{
		Version:      "v1",
		Addr:         ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
}

type Server struct {
	*http.Server
}

func NewServer(cfg *Config, ctrls ...Controller) *Server {
	return &Server{
		Server: &http.Server{
			Addr:         cfg.Addr,
			Handler:      newRouter(cfg.Version, ctrls...),
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
		},
	}
}

func (s *Server) Run() error {
	return s.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}
