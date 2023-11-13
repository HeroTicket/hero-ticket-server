package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/heroticket/internal/ws"
	"github.com/heroticket/pkg/mongo"
	"github.com/heroticket/pkg/shutdown"
	auth "github.com/iden3/go-iden3-auth/v2"
	"github.com/iden3/go-iden3-auth/v2/loaders"
	"github.com/iden3/go-iden3-auth/v2/pubsignals"
	"github.com/iden3/go-iden3-auth/v2/state"
	"github.com/iden3/iden3comm/v2/protocol"
	"go.uber.org/zap"
	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"

	_ "github.com/joho/godotenv/autoload"
)

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

	zap.L().Info("connected to mongo")

	h := newHandler()

	go func() {
		zap.L().Info("starting ngrok server")

		// if err := runNgrok(h); err != nil && err != http.ErrServerClosed {
		// 	panic(err)
		// }

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

var requestMap = make(map[string]interface{})

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

	var request protocol.AuthorizationRequestMessage = auth.CreateAuthorizationRequestWithMessage(
		"Login to Hero Ticket",
		"Scan the QR code to login to Hero Ticket",
		audience,
		uri,
	)

	request.ID = sessionId
	request.ThreadID = sessionId

	requestMap[sessionId] = request

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

	authRequest, ok := requestMap[sessionId]
	if !ok {
		http.Error(w, "request not found", http.StatusNotFound)
		return
	}

	go ws.Send(ws.Message{
		ID:   id,
		Type: ws.EventMessage,
		Event: ws.Event{
			Name:   "login-callback",
			Status: ws.InProgress,
		},
	})

	ipfsUrl := "https://ipfs.io"
	contractAddress := "134B1BE34911E39A8397ec6289782989729807a4"
	resolverPrefix := "polygon:mumbai"
	ketDir := "./keys"

	var verificationKeyLoader = &loaders.FSKeyLoader{
		Dir: ketDir,
	}

	resolver := state.ETHResolver{
		RPCUrl:          os.Getenv("RPC_URL_MUMBAI"),
		ContractAddress: common.HexToAddress(contractAddress),
	}

	resolvers := map[string]pubsignals.StateResolver{
		resolverPrefix: &resolver,
	}

	verifier, err := auth.NewVerifier(
		verificationKeyLoader,
		resolvers,
		auth.WithIPFSGateway(ipfsUrl),
	)
	if err != nil {
		http.Error(w, "could not create verifier", http.StatusInternalServerError)
		return
	}

	authResponse, err := verifier.FullVerify(
		r.Context(),
		string(tokenBytes),
		authRequest.(protocol.AuthorizationRequestMessage),
		pubsignals.WithAcceptedProofGenerationDelay(time.Minute*5),
	)
	if err != nil {
		http.Error(w, "could not verify", http.StatusInternalServerError)
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
