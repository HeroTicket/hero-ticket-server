package rest

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/heroticket/internal/db"
	"github.com/heroticket/internal/service/auth"
	"github.com/heroticket/internal/service/did"
	"github.com/heroticket/internal/service/jwt"
	"github.com/heroticket/internal/service/user"
	"github.com/heroticket/internal/web3"
	"github.com/heroticket/internal/ws"
	"go.uber.org/zap"
)

type UserCtrl struct {
	baseUrl string

	auth auth.Service
	did  did.Service
	jwt  jwt.Service
	user user.Service
	tx   db.Tx
}

func NewUserCtrl(auth auth.Service, did did.Service, jwt jwt.Service, user user.Service, tx db.Tx, baseUrl string) *UserCtrl {
	return &UserCtrl{
		baseUrl: baseUrl,
		auth:    auth,
		did:     did,
		jwt:     jwt,
		user:    user,
		tx:      tx,
	}
}

func (c *UserCtrl) Pattern() string {
	return "/users"
}

func (c *UserCtrl) Handler() http.Handler {
	r := chi.NewRouter()

	r.Get("/login-qr", c.loginQR)
	r.Post("/login-callback", c.loginCallback)

	r.Group(func(r chi.Router) {
		r.Use(TokenRequired(c.jwt))
		r.Post("/register", c.register)
	})

	// r.Group(func(r chi.Router) {
	// 	r.Get("/profile", c.profile)
	// 	r.Get("/purchased-tickets", c.purchasedTickets)
	// 	r.Get("/issued-tickets", c.issuedTickets)
	// })

	return r
}

// LoginQR godoc
//
//	@Tags			users
//	@Summary		returns login qr code
//	@Description	returns login qr code
//	@Produce		json
//	@Param			sessionId	query		string	true	"session id"
//	@Success		200			{object}	CommonResponse
//	@Failure		400			{object}	CommonResponse
//	@Failure		500			{object}	CommonResponse
//	@Router			/users/login-qr [get]
func (c *UserCtrl) loginQR(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// 1. get session id from query params
	sessionId := r.URL.Query().Get("sessionId")

	// 2. validate session id
	id := ws.ID(sessionId)

	if !id.Valid() {
		ErrorJSON(w, "invalid session id")
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

	callbackUrl := fmt.Sprintf("%s/v1/users/login-callback?sessionId=%s", c.baseUrl, sessionId)

	// TODO: fetch audience from db
	audience := os.Getenv("VERIFIER_DID")

	// 3. create login request
	req, err := c.auth.AuthorizationRequest(r.Context(), auth.AuthorizationRequestParams{
		ID:          sessionId,
		Reason:      "Login to Hero Ticket",
		Message:     "Scan the QR code to login to Hero Ticket",
		CallbackUrl: callbackUrl,
		Sender:      audience,
	})
	if err != nil {
		zap.L().Error("failed to create login request", zap.Error(err))
		ErrorJSON(w, "failed to create login request", http.StatusInternalServerError)
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

	resp := CommonResponse{
		Status:  http.StatusOK,
		Message: "Successfully created login request",
		Data:    req,
	}

	// 4. return login request as json response
	_ = WriteJSON(w, http.StatusOK, resp)
}

// LoginCallback godoc
//
//	@Tags			users
//	@Summary		processes login callback
//	@Description	processes login callback
//	@Accept 		plain
//	@Produce		json
//	@Param			sessionId	query		string	true	"session id"
//	@Success		200			{object}	CommonResponse
//	@Failure		400			{object}	CommonResponse
//	@Failure		500			{object}	CommonResponse
//	@Router			/users/login-callback [post]
func (c *UserCtrl) loginCallback(w http.ResponseWriter, r *http.Request) {
	// 1. get session id from query params
	sessionId := r.URL.Query().Get("sessionId")

	// 2. validate session id
	id := ws.ID(sessionId)

	if !id.Valid() {
		ErrorJSON(w, "invalid session id")
		return
	}

	// 3. read token from request body
	tokenBytes, err := io.ReadAll(r.Body)
	if err != nil {
		zap.L().Error("failed to read token", zap.Error(err))
		ErrorJSON(w, "failed to read token", http.StatusInternalServerError)
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

	// 4. handle login callback
	resp, err := c.auth.AuthorizationCallback(r.Context(), sessionId, string(tokenBytes))
	if err != nil {
		zap.L().Error("failed to handle login callback", zap.Error(err))
		ErrorJSON(w, "failed to handle login callback", http.StatusInternalServerError)
		return
	}

	userIdentifier := resp.From

	// 5. generate jwt token
	token, err := c.jwt.GenerateToken(jwt.JWTUser{
		Identifier: userIdentifier,
	})
	if err != nil {
		zap.L().Error("failed to generate jwt token", zap.Error(err))
		ErrorJSON(w, "failed to generate jwt token", http.StatusInternalServerError)
		return
	}

	go ws.Send(ws.Message{
		ID:   id,
		Type: ws.EventMessage,
		Event: ws.Event{
			Name:   "login-callback",
			Status: ws.Done,
			Data:   token,
		},
	})

	// 6. return success response
	response := CommonResponse{
		Status:  http.StatusOK,
		Message: fmt.Sprintf("User with ID %s Successfully authenticated", userIdentifier),
	}

	_ = WriteJSON(w, http.StatusOK, response)
}

// Register godoc
//
//	@Tags			users
//	@Summary		registers user
//	@Description	registers user
//	@Produce		json
//	@Param			sessionId		query		string	true	"session id"
//	@Param			walletAddress	query		string	true	"wallet address"
//	@Success		201		{object}	CommonResponse
//	@Failure		400		{object}	CommonResponse
//	@Failure		500		{object}	CommonResponse
//	@Router			/users/register [post]
func (c *UserCtrl) register(w http.ResponseWriter, r *http.Request) {
	// 1. get user from context
	jwtUser, err := c.jwt.FromContext(r.Context())
	if err != nil {
		zap.L().Error("failed to get user from context", zap.Error(err))
		ErrorJSON(w, "user not found")
		return
	}

	// 2. get session id from query params
	sessionId := r.URL.Query().Get("sessionId")

	// 3. validate session id
	id := ws.ID(sessionId)

	if !id.Valid() {
		ErrorJSON(w, "invalid session id")
		return
	}

	// 4. get wallet address from query params
	walletAddress := r.URL.Query().Get("walletAddress")

	// 5. validate wallet address
	if !web3.IsAddressValid(walletAddress) {
		ErrorJSON(w, "invalid wallet address", http.StatusInternalServerError)
		return
	}

	// 7. check if identifier is already registered or not
	u, err := c.user.FindUserByID(r.Context(), jwtUser.Identifier)
	if err != nil {
		if err != user.ErrUserNotFound {
			zap.L().Error("failed to find user", zap.Error(err))
			ErrorJSON(w, "failed to find user", http.StatusInternalServerError)
			return
		}
	} else {
		if u.AccountAddress != "" {
			ErrorJSON(w, "user already registered", http.StatusInternalServerError)
			return
		}
	}

	// 8. create user tba

	// 9. update user

	// 10. return success response
}

// Profile godoc
//
//	@Tags			users
//	@Summary		returns user profile
//	@Description	returns user profile
//	@Produce		json
//	@Success		200			{object}	CommonResponse
//	@Failure		400			{object}	CommonResponse
//	@Router			/users/profile [get]
// func (c *UserCtrl) profile(w http.ResponseWriter, r *http.Request) {
// 	// 1. get user from context
// 	jwtUser, err := c.jwt.FromContext(r.Context())
// 	if err != nil {
// 		zap.L().Error("failed to get user from context", zap.Error(err))
// 		ErrorJSON(w, "user not found")
// 		return
// 	}

// 	// 2. get user from db
// 	u, err := c.user.FindUserByDID(r.Context(), jwtUser.DID)
// 	if err != nil {
// 		zap.L().Error("failed to find user", zap.Error(err))
// 		ErrorJSON(w, "failed to find user")
// 		return
// 	}

// 	resp := CommonResponse{
// 		Status:  http.StatusOK,
// 		Message: "Successfully fetched user profile",
// 		Data:    u,
// 	}

// 	// 3. return user as json response
// 	_ = WriteJSON(w, http.StatusOK, resp)
// }

// PurchasedTickets godoc
//
//	@Tags			users
//	@Summary		returns purchased tickets
//	@Description	returns purchased tickets
//	@Produce		json

// @Router			/users/purchased-tickets [get]
// func (c *UserCtrl) purchasedTickets(w http.ResponseWriter, r *http.Request) {
// 	// 1. get user from context
// 	_, err := c.jwt.FromContext(r.Context())
// 	if err != nil {
// 		zap.L().Error("failed to get user from context", zap.Error(err))
// 		ErrorJSON(w, "user not found")
// 		return
// 	}

// 	// 2. get purchased tickets from db

// 	// 3. return purchased tickets as json response
// }

// IssuedTickets godoc
//
//	@Tags			users
//	@Summary		returns issued tickets
//	@Description	returns issued tickets
//	@Produce		json

// @Router			/users/issued-tickets [get]
// func (c *UserCtrl) issuedTickets(w http.ResponseWriter, r *http.Request) {
// 	// 1. get user from context
// 	_, err := c.jwt.FromContext(r.Context())
// 	if err != nil {
// 		zap.L().Error("failed to get user from context", zap.Error(err))
// 		ErrorJSON(w, "user not found")
// 		return
// 	}

// 	// 2. get issued tickets from db

// 	// 3. return issued tickets as json response
// }
