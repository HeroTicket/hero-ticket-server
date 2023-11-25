package rest

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/heroticket/internal/db"
	"github.com/heroticket/internal/service/auth"
	"github.com/heroticket/internal/service/jwt"
	"github.com/heroticket/internal/service/user"
	"github.com/heroticket/internal/web3"
	"github.com/heroticket/internal/ws"
	"go.uber.org/zap"
)

type UserCtrl struct {
	baseUrl string

	auth auth.Service
	jwt  jwt.Service
	user user.Service
	tx   db.Tx
}

func NewUserCtrl(auth auth.Service, jwt jwt.Service, user user.Service, tx db.Tx, baseUrl string) *UserCtrl {
	return &UserCtrl{
		baseUrl: baseUrl,
		auth:    auth,
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
	r.Post("/refresh", c.refresh)

	r.Group(func(r chi.Router) {
		r.Use(TokenRequired(c.jwt))
		r.Post("/register/{accountAddress}", c.register)
	})

	return r
}

// LoginQR godoc
//
//	@Tags			users
//	@Summary		returns login qr code
//	@Description	returns login qr code
//	@Accept 		json
//	@Produce		json
//	@Param			sessionId	query		string	true	"session id"
//	@Success		200			{object}	CommonResponse{data=protocol.AuthorizationRequestMessage}
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
//	@Param			token		body		string	true	"token"
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

	userID := resp.From

	// 5. generate jwt token
	tokenPair, err := c.jwt.GenerateTokenPair(jwt.JWTUser{
		ID: userID,
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
			Data:   tokenPair,
		},
	})

	// 6. return success response
	response := CommonResponse{
		Status:  http.StatusOK,
		Message: fmt.Sprintf("User with ID %s Successfully authenticated", userID),
	}

	_ = WriteJSON(w, http.StatusOK, response)
}

// Refresh Token Pair godoc
//
//	@Tags			users
//	@Summary		refreshes token pair
//	@Description	refreshes token pair
//	@Accept 		plain
//	@Produce		json
//	@Param			refreshToken		body		string	true	"refresh token"
//	@Success		200			{object}	CommonResponse{data=jwt.TokenPair}
//	@Failure		400			{object}	CommonResponse
//	@Failure		500			{object}	CommonResponse
//	@Router			/users/refresh [post]
func (c *UserCtrl) refresh(w http.ResponseWriter, r *http.Request) {
	// 1. read token from request body
	tokenBytes, err := io.ReadAll(r.Body)
	if err != nil {
		zap.L().Error("failed to read token", zap.Error(err))
		ErrorJSON(w, "failed to read token", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// 2. validate token
	jwtUser, err := c.jwt.VerifyToken(string(tokenBytes), jwt.TokenRoleRefresh)
	if err != nil {
		ErrorJSON(w, "invalid token", http.StatusBadRequest)
		return
	}

	// 3. generate new token pair
	newTokenPair, err := c.jwt.GenerateTokenPair(jwt.JWTUser{
		ID: jwtUser.ID,
	})
	if err != nil {
		zap.L().Error("failed to generate token pair", zap.Error(err))
		ErrorJSON(w, "failed to generate token pair", http.StatusInternalServerError)
		return
	}

	// 4. return new token pair as json response
	resp := CommonResponse{
		Status:  http.StatusOK,
		Message: "Successfully refreshed token pair",
		Data:    newTokenPair,
	}

	_ = WriteJSON(w, http.StatusOK, resp)
}

// Register godoc
//
//	@Tags			users
//	@Summary		registers user
//	@Description	registers user
//	@Produce		json
//	@Param			accountAddress 	path	string	true	"account address"
//	@Success		201		{object}	CommonResponse
//	@Failure		400		{object}	CommonResponse
//	@Failure		500		{object}	CommonResponse
//	@Security 		BearerAuth
//	@Router			/users/register/{accountAddress} [post]
func (c *UserCtrl) register(w http.ResponseWriter, r *http.Request) {
	// 1. get user from context
	jwtUser, err := c.jwt.FromContext(r.Context())
	if err != nil {
		zap.L().Error("failed to get user from context", zap.Error(err))
		ErrorJSON(w, "user not found")
		return
	}

	// 2. get account address from path params
	accountAddress := chi.URLParam(r, "accountAddress")

	// 3. validate account address
	if !web3.IsAddressValid(accountAddress) {
		ErrorJSON(w, "invalid account address", http.StatusInternalServerError)
		return
	}

	// 4. check if id is already registered or not
	_, err = c.user.FindUserByID(r.Context(), jwtUser.ID)
	if err == nil {
		ErrorJSON(w, "user already registered", http.StatusBadRequest)
		return
	}

	if err != nil && err != user.ErrUserNotFound {
		zap.L().Error("failed to find user", zap.Error(err))
		ErrorJSON(w, "failed to register user", http.StatusInternalServerError)
		return
	}

	// 5. create user tba

	// TODO: call contract to create user tba

	tempTbaAddress := "0x1234567890"

	// 6. create user
	params := user.CreateUserParams{
		ID:             jwtUser.ID,
		AccountAddress: accountAddress,
		TbaAddress:     tempTbaAddress,
		Name:           accountAddress,
		IsAdmin:        false,
	}

	u, err := c.user.CreateUser(r.Context(), params)
	if err != nil {
		zap.L().Error("failed to create user", zap.Error(err))
		ErrorJSON(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	// 7. return user as json response
	resp := CommonResponse{
		Status:  http.StatusCreated,
		Message: "Successfully created user",
		Data:    u,
	}

	_ = WriteJSON(w, http.StatusCreated, resp)
}
