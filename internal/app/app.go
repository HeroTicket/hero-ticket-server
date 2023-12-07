package app

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
	IdleTimeout  time.Duration
}

func DefaultConfig() *Config {
	return &Config{
		Version:      "v1",
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  3 * time.Minute,
	}
}

type App struct {
	*http.Server
}

func New(cfg *Config, ctrls ...Controller) *App {
	return &App{
		Server: &http.Server{
			Addr:         cfg.Addr,
			Handler:      newRouter(cfg.Version, ctrls...),
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
	}
}

func (s *App) Run() error {
	return s.ListenAndServe()
}

func (s *App) Shutdown(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}
