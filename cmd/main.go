package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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

	h := newHandler()

	go func() {
		zap.L().Info("starting ngrok server")

		if err := runNgrok(h); err != nil && err != http.ErrServerClosed {
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

func newHandler() http.Handler {
	r := chi.NewRouter()

	r.Use(cors.AllowAll().Handler)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/users", func(r chi.Router) {
				r.Get("/login-qr", nil)          // get qr code
				r.Post("/login-callback", nil)   // callback from qr code
				r.Post("/logout", nil)           // logout user
				r.Post("/refresh-token", nil)    // refresh token
				r.Post("/create-tba", nil)       // create token bound account (tba) for user
				r.Get("/purchased-tickets", nil) // get list of purchased tickets
				r.Get("/issued-tickets", nil)    // get list of published tickets
			})

			r.Route("/tickets", func(r chi.Router) {
				r.Post("/create", nil)                // create ticket
				r.Get("/create-qr", nil)              // get qr code
				r.Post("/create-callback", nil)       // callback from qr code
				r.Get("/list", nil)                   // get list of tickets
				r.Get("/list/{id}", nil)              // get ticket by id
				r.Get("{id}/purchase-qr", nil)        // get qr code
				r.Post("{id}/purchase-callback", nil) // callback from qr code
				r.Get("/{id}/verify-qr", nil)         // get qr code
				r.Post("/{id}/verify-callback", nil)  // callback from qr code
			})

			r.HandleFunc("/ws", nil) //	handle websocket
		})
	})

	return r
}
