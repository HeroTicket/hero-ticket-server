package rest

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/heroticket/internal/app/ws"
	"github.com/heroticket/internal/logger"
	"github.com/heroticket/internal/service/auth"
	"github.com/heroticket/internal/service/jwt"
	"github.com/heroticket/internal/service/ticket"
	"github.com/heroticket/internal/service/user"
	"github.com/heroticket/internal/web3"
)

type UserCtrl struct {
	serverUrl string

	auth   auth.Service
	jwt    jwt.Service
	user   user.Service
	ticket ticket.Service
}

func NewUserCtrl(auth auth.Service, jwt jwt.Service, user user.Service, ticket ticket.Service, serverUrl string) *UserCtrl {
	return &UserCtrl{
		serverUrl: serverUrl,
		auth:      auth,
		jwt:       jwt,
		user:      user,
		ticket:    ticket,
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
		// 사용자 정보 조회
		r.Get("/info", c.info)
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
//	@Router			/v1/users/login-qr [get]
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

	admin, err := c.user.FindAdmin(r.Context())
	if err != nil {
		logger.Error("failed to find admin", "error", err)
		ErrorJSON(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	callbackUrl := fmt.Sprintf("%s/v1/users/login-callback?sessionId=%s", c.serverUrl, sessionId)

	// 3. create login request
	req, err := c.auth.AuthorizationRequest(r.Context(), auth.AuthorizationRequestParams{
		ID:          sessionId,
		Reason:      "Login to Hero Ticket",
		Message:     "Scan the QR code to login to Hero Ticket",
		CallbackUrl: callbackUrl,
		Sender:      admin.ID,
	})
	if err != nil {
		logger.Error("failed to create login request", "error", err)
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
//	@Router			/v1/users/login-callback [post]
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
		logger.Error("failed to read token", "error", err)
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
		logger.Error("failed to handle login callback", "error", err)
		ErrorJSON(w, "failed to handle login callback", http.StatusInternalServerError)
		return
	}

	userID := resp.From

	// 5. generate jwt token
	tokenPair, err := c.jwt.GenerateTokenPair(jwt.JWTUser{
		ID: userID,
	})
	if err != nil {
		logger.Error("failed to generate jwt token", "error", err)
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
//	@Router			/v1/users/refresh [post]
func (c *UserCtrl) refresh(w http.ResponseWriter, r *http.Request) {
	// 1. read token from request body
	tokenBytes, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error("failed to read token", "error", err)
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
		logger.Error("failed to generate token pair", "error", err)
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

// Info godoc
//
//	@Tags			users
//	@Summary		returns user info
//	@Description	returns user info
//	@Accept			json
//	@Produce		json
//	@Success		200			{object}	CommonResponse{data=user.User}
//	@Failure		400			{object}	CommonResponse
//	@Failure		404			{object}	CommonResponse
//	@Failure		500			{object}	CommonResponse
//	@Security 		BearerAuth
//	@Router			/v1/users/info [get]
func (c *UserCtrl) info(w http.ResponseWriter, r *http.Request) {
	// 1. get user from context
	jwtUser, err := c.jwt.FromContext(r.Context())
	if err != nil {
		logger.Error("failed to get user from context", "error", err)
		ErrorJSON(w, "user not found")
		return
	}

	// 2. get user from db
	u, err := c.user.FindUserByID(r.Context(), jwtUser.ID)
	if err != nil {
		if err == user.ErrUserNotFound {
			ErrorJSON(w, "user not registered yet", http.StatusNotFound)
			return
		}
		logger.Error("failed to find user", "error", err)
		ErrorJSON(w, "failed to find user", http.StatusInternalServerError)
		return
	}

	// 3. return user as json response
	resp := CommonResponse{
		Status:  http.StatusOK,
		Message: "Successfully fetched user info",
		Data:    u,
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
//	@Success		201		{object}	CommonResponse{data=user.User}
//	@Failure		400		{object}	CommonResponse
//	@Failure		500		{object}	CommonResponse
//	@Security 		BearerAuth
//	@Router			/v1/users/register/{accountAddress} [post]
func (c *UserCtrl) register(w http.ResponseWriter, r *http.Request) {
	// 1. get user from context
	jwtUser, err := c.jwt.FromContext(r.Context())
	if err != nil {
		logger.Error("failed to get user from context", "error", err)
		ErrorJSON(w, "user not found")
		return
	}

	// 2. get account address from path params
	rawAccountAddress := strings.ToLower(chi.URLParam(r, "accountAddress"))

	// 3. validate account address
	if !web3.IsAddressValid(rawAccountAddress) {
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
		ErrorJSON(w, "failed to check if user is already registered or not", http.StatusInternalServerError)
		return
	}

	accountAddress := web3.HexToAddress(rawAccountAddress)

	// TODO: uri should be dynamic
	uri := "https://gold-cool-goat-213.mypinata.cloud/ipfs/QmfFbvLH37DebBqmVBm7V8ecfzgjFPnPeHRYiYk1PNoW84/2level.png"

	// 5. create user tba
	tbaCreated, err := c.ticket.CreateTBA(r.Context(), accountAddress, uri)
	if err != nil {
		logger.Error("failed to create tba", "error", err)
		ErrorJSON(w, "failed to create tba", http.StatusInternalServerError)
		return
	}

	tbaAddress := strings.ToLower(tbaCreated.Account.Hex())

	// 6. create user
	params := user.CreateUserParams{
		ID:             jwtUser.ID,
		AccountAddress: rawAccountAddress,
		TbaAddress:     tbaAddress,
		Name:           rawAccountAddress,
		IsAdmin:        false,
	}

	u, err := c.user.CreateUser(r.Context(), params)
	if err != nil {
		logger.Error("failed to create user", "error", err)
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
