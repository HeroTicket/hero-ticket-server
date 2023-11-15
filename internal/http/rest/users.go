package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/heroticket/internal/db"
	"github.com/heroticket/internal/did"
	"github.com/heroticket/internal/jwt"
	"github.com/heroticket/internal/user"
	"github.com/heroticket/internal/ws"
)

type UserCtrl struct {
	did  did.Service
	jwt  jwt.Service
	user user.Service
	tx   db.Tx
}

func NewUserCtrl(did did.Service, jwt jwt.Service, user user.Service, tx db.Tx) *UserCtrl {
	return &UserCtrl{
		did:  did,
		jwt:  jwt,
		user: user,
		tx:   tx,
	}
}

func (c *UserCtrl) Pattern() string {
	return "/users"
}

func (c *UserCtrl) Handler() http.Handler {
	r := chi.NewRouter()

	r.Get("/login-qr", c.loginQR)
	r.Post("/login-callback", c.loginCallback)
	r.Post("/logout", nil)
	r.Post("/refresh-token", nil)
	r.Post("/create-tba", nil)
	r.Get("/purchased-tickets", nil)
	r.Get("/issued-tickets", nil)

	return r
}

func (c *UserCtrl) loginQR(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	sessionId := r.URL.Query().Get("sessionId")

	id := ws.ID(sessionId)

	if !id.Valid() {
		http.Error(w, "invalid session id", http.StatusBadRequest)
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

	callbackUrl := fmt.Sprintf("%s/api/v1/users/login-callback?sessionId=%s", os.Getenv("NGROK_URL"), sessionId)

	audience := os.Getenv("VERIFIER_DID")

	request, err := c.did.LoginRequest(r.Context(), sessionId, audience, callbackUrl)
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

	response, err := json.Marshal(request)
	if err != nil {
		http.Error(w, "failed to marshal login request", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func (c *UserCtrl) loginCallback(w http.ResponseWriter, r *http.Request) {
	sessionId := r.URL.Query().Get("sessionId")

	id := ws.ID(sessionId)

	if !id.Valid() {
		http.Error(w, "invalid session id", http.StatusBadRequest)
		return
	}

	tokenBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read token", http.StatusInternalServerError)
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

	authResponse, err := c.did.LoginCallback(r.Context(), sessionId, string(tokenBytes))
	if err != nil {
		http.Error(w, "failed to login", http.StatusInternalServerError)
		return
	}

	go ws.Send(ws.Message{
		ID:   id,
		Type: ws.EventMessage,
		Event: ws.Event{
			Name:   "login-callback",
			Status: ws.Done,
		},
	})

	userDID := authResponse.From

	messageBytes := []byte("User with ID " + userDID + " Successfully authenticated")

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(messageBytes)
}
