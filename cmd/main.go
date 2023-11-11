package main

import (
	"context"
	"log"
	"net/http"
	"syscall"

	"github.com/heroticket/pkg/shutdown"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("recovered from panic: %v", r)
		}
	}()

	srv := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("hello world"))
		}),
	}

	go func() {
		log.Println("starting server...")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	stop := shutdown.GracefulShutdown(func() {
		log.Println("shutting down server...")
		if err := srv.Shutdown(context.Background()); err != nil {
			panic(err)
		}

		log.Println("server gracefully stopped")
	}, syscall.SIGINT, syscall.SIGTERM)

	<-stop
}
