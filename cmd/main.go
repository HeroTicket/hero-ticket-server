package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"syscall"

	"github.com/heroticket/pkg/mongo"
	"github.com/heroticket/pkg/shutdown"
	"go.uber.org/zap"
	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
)

func main() {
	logger, _ := zap.NewProduction(zap.Fields(zap.String("service", "hero-ticket")))
	defer logger.Sync()

	zap.ReplaceGlobals(logger)

	defer func() {
		if r := recover(); r != nil {
			zap.L().Info("recovered from panic", zap.Any("r", r))
		}
	}()

	client, err := mongo.New(context.Background(), "mongodb://root:example@localhost:27017/")
	if err != nil {
		panic(err)
	}

	zap.L().Info("connected to mongo")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello world"))
	})

	go func() {
		zap.L().Info("starting ngrok server")

		if err := runNgrok(handler); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	stop := shutdown.GracefulShutdown(func() {
		if err := client.Disconnect(context.Background()); err != nil {
			panic(err)
		}

		zap.L().Info("disconnected from mongo")
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

	if err := os.Setenv("NGROK_URL", url); err != nil {
		return err
	}

	return http.Serve(tun, h)
}
