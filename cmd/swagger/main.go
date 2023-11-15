package main

import (
	"context"
	"net/http"
	"syscall"

	"github.com/go-chi/chi/v5"
	_ "github.com/heroticket/docs"
	"github.com/heroticket/internal/shutdown"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	r := chi.NewRouter()

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:1323/swagger/doc.json"),
	))

	srv := &http.Server{
		Addr:    ":1323",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// http://localhost:1323/swagger/index.html
	logger.Info("Browse to http://localhost:1323/swagger/index.html")

	<-shutdown.GracefulShutdown(func() {
		if err := srv.Shutdown(context.Background()); err != nil {
			panic(err)
		}
	}, syscall.SIGINT, syscall.SIGTERM)
}
