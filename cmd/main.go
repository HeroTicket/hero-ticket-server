package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"syscall"

	"github.com/heroticket/pkg/mongo"
	"github.com/heroticket/pkg/shutdown"
	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("recovered from panic: %v", r)
		}
	}()

	client, err := mongo.New(context.Background(), "mongodb://root:example@localhost:27017/")
	if err != nil {
		panic(err)
	}

	log.Println("connected to mongo")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello world"))
	})

	go func() {
		log.Println("starting server...")

		if err := runNgrok(handler); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	stop := shutdown.GracefulShutdown(func() {
		if err := client.Disconnect(context.Background()); err != nil {
			panic(err)
		}

		log.Println("disconnected from mongo")
	}, syscall.SIGINT, syscall.SIGTERM)

	<-stop
}

func runNgrok(h http.Handler) error {
	tun, err := ngrok.Listen(
		context.Background(),
		config.HTTPEndpoint(),
	)
	if err != nil {
		return err
	}

	url := tun.URL()

	fmt.Println("ngrok url: ", url)

	return http.Serve(tun, h)
}
