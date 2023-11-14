package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/heroticket/internal/did"
	dredis "github.com/heroticket/internal/did/cache/redis"
	"github.com/heroticket/internal/infra/mongo"
	"github.com/heroticket/internal/infra/redis"
	"github.com/heroticket/internal/ws"
	"github.com/heroticket/pkg/shutdown"
	"go.uber.org/zap"

	_ "github.com/joho/godotenv/autoload"
)

var didService did.Service

func main() {
	logger, _ := zap.NewProduction(zap.Fields(zap.String("service", "hero-ticket")))
	defer logger.Sync()

	zap.ReplaceGlobals(logger)

	zap.L().Info("starting server")

	defer func() {
		if r := recover(); r != nil {
			zap.L().Info("recovered from panic", zap.Any("r", r))
		}
	}()

	client, err := mongo.New(context.Background(), "mongodb://root:example@localhost:27017/")
	if err != nil {
		panic(err)
	}

	cache, err := redis.NewCache(context.Background(), "localhost:6379")
	if err != nil {
		panic(err)
	}

	zap.L().Info("connected to redis")

	didService = did.New(dredis.New(cache), os.Getenv("RPC_URL_MUMBAI"))

	zap.L().Info("connected to mongo")

	h := newHandler()

	go func() {
		zap.L().Info("starting server")
		if err := runServer(h); err != nil && err != http.ErrServerClosed {
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

func runServer(h http.Handler) error {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: h,
	}

	return srv.ListenAndServe()
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
				r.Get("/login-qr", loginQR)              // get qr code
				r.Post("/login-callback", loginCallback) // callback from qr code
				r.Post("/logout", nil)                   // logout user
				r.Post("/refresh-token", nil)            // refresh token
				r.Post("/create-tba", nil)               // create token bound account (tba) for user
				r.Get("/purchased-tickets", nil)         // get list of purchased tickets
				r.Get("/issued-tickets", nil)            // get list of published tickets
			})

			r.Route("/tickets", func(r chi.Router) {
				r.Post("/create", nil)                 // create ticket
				r.Get("/create-qr", nil)               // get qr code
				r.Post("/create-callback", nil)        // callback from qr code
				r.Get("/list", nil)                    // get list of tickets
				r.Get("/list/{id}", nil)               // get ticket by id
				r.Get("/{id}/purchase-qr", nil)        // get qr code
				r.Post("/{id}/purchase-callback", nil) // callback from qr code
				r.Get("/{id}/verify-qr", nil)          // get qr code
				r.Post("/{id}/verify-callback", nil)   // callback from qr code
			})
		})
	})

	r.HandleFunc("/ws", ws.Serve()) //	handle websocket connections

	return r
}

func loginQR(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	sessionId := r.URL.Query().Get("sessionId")

	id := ws.ID(sessionId)

	if !id.Valid() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	go ws.Send(ws.Message{
		ID:   id,
		Type: ws.EventMessage,
		Event: ws.Event{
			Name:   "login-qr",
			Status: ws.InProgress,
		},
	})

	uri := fmt.Sprintf("%s/api/v1/users/login-callback?sessionId=%s", os.Getenv("NGROK_URL"), sessionId)

	audience := os.Getenv("VERIFIER_DID")

	request, err := didService.LoginRequest(r.Context(), sessionId, audience, uri)
	if err != nil {
		http.Error(w, "failed to create login request", http.StatusInternalServerError)
		return
	}

	go ws.Send(ws.Message{
		ID:   id,
		Type: ws.EventMessage,
		Event: ws.Event{
			Name:   "login-qr",
			Status: ws.Done,
		},
	})

	resp, err := json.Marshal(request)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func loginCallback(w http.ResponseWriter, r *http.Request) {
	sessionId := r.URL.Query().Get("sessionId")

	id := ws.ID(sessionId)

	if !id.Valid() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokenBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	go ws.Send(ws.Message{
		ID:   id,
		Type: ws.EventMessage,
		Event: ws.Event{
			Name:   "login-callback",
			Status: ws.InProgress,
		},
	})

	authResponse, err := didService.LoginCallback(r.Context(), sessionId, string(tokenBytes))
	if err != nil {
		http.Error(w, "failed to login", http.StatusInternalServerError)
		return
	}

	// TODO: generate jwt token

	go ws.Send(ws.Message{
		ID:   id,
		Type: ws.EventMessage,
		Event: ws.Event{
			Name:   "login-callback",
			Status: ws.Done,
		},
	})

	userID := authResponse.From

	messageBytes := []byte("User with ID " + userID + " Successfully authenticated")

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(messageBytes)
}
