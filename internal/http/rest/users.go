package rest

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/anandvarma/namegen"
	"github.com/go-chi/chi/v5"
	"github.com/heroticket/internal/db"
	"github.com/heroticket/internal/did"
	"github.com/heroticket/internal/jwt"
	"github.com/heroticket/internal/user"
	"github.com/heroticket/internal/ws"
	"go.uber.org/zap"
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
	r.Post("/refresh-token", c.refreshToken)

	r.Group(func(r chi.Router) {
		r.Use(AccessTokenRequired(c.jwt))

		r.Post("/create-tba", c.createTBA)
		r.Get("/profile", c.profile)
		r.Get("/purchased-tickets", c.purchasedTickets)
		r.Get("/issued-tickets", c.issuedTickets)
	})

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

	callbackUrl := fmt.Sprintf("%s/api/v1/users/login-callback?sessionId=%s", os.Getenv("NGROK_URL"), sessionId)

	audience := os.Getenv("VERIFIER_DID")

	// 3. create login request
	request, err := c.did.LoginRequest(r.Context(), sessionId, audience, callbackUrl)
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
		Data:    request,
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
	authResponse, err := c.did.AuthorizationCallback(r.Context(), sessionId, string(tokenBytes))
	if err != nil {
		zap.L().Error("failed to handle login callback", zap.Error(err))
		ErrorJSON(w, "failed to handle login callback", http.StatusInternalServerError)
		return
	}

	userDID := authResponse.From

	// 5. find or create user
	u, err := c.user.FindUserByDID(r.Context(), userDID)
	if err != nil {
		if err == user.ErrUserNotFound {
			u, err = c.user.CreateUser(r.Context(), &user.User{
				DID:  userDID,
				Name: namegen.New().Get(),
			})
			if err != nil {
				zap.L().Error("failed to create user", zap.Error(err))
				ErrorJSON(w, "failed to create user", http.StatusInternalServerError)
				return
			}
		} else {
			zap.L().Error("failed to find user", zap.Error(err))
			ErrorJSON(w, "failed to find user", http.StatusInternalServerError)
			return
		}
	}

	// 6. generate token pair
	jwtUser := jwt.JWTUser{
		DID: u.DID,
	}

	tokenPair, err := c.jwt.GenerateTokenPair(jwtUser)
	if err != nil {
		zap.L().Error("failed to generate token pair", zap.Error(err))
		ErrorJSON(w, "failed to generate token pair", http.StatusInternalServerError)
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

	// 7. return success response
	response := CommonResponse{
		Status:  http.StatusOK,
		Message: fmt.Sprintf("User with ID %s Successfully authenticated", userDID),
	}

	_ = WriteJSON(w, http.StatusOK, response)
}

// RefreshToken godoc
//
//	@Tags			users
//	@Summary		refreshes token pair
//	@Description	refreshes token pair
//	@Accept 		plain
//	@Produce		json
//	@Success		200			{object}	CommonResponse
//	@Failure		400			{object}	CommonResponse
//	@Router			/users/refresh-token [post]
func (c *UserCtrl) refreshToken(w http.ResponseWriter, r *http.Request) {
	// 1. get refresh token from request body
	refreshToken, err := io.ReadAll(r.Body)
	if err != nil {
		zap.L().Error("failed to read refresh token", zap.Error(err))
		ErrorJSON(w, "failed to read refresh token")
		return
	}
	defer r.Body.Close()

	// 2. validate refresh token
	jwtUser, err := c.jwt.VerifyToken(string(refreshToken), jwt.TokenRoleRefresh)
	if err != nil {
		zap.L().Error("failed to validate refresh token", zap.Error(err))
		ErrorJSON(w, "failed to validate refresh token")
		return
	}

	// 3. get user from db
	u, err := c.user.FindUserByDID(r.Context(), jwtUser.DID)
	if err != nil {
		zap.L().Error("failed to find user", zap.Error(err))
		ErrorJSON(w, "failed to find user")
		return
	}

	newJwtUser := jwt.JWTUser{
		DID: u.DID,
	}

	// 4. generate new token pair
	tokenPair, err := c.jwt.GenerateTokenPair(newJwtUser)
	if err != nil {
		zap.L().Error("failed to generate token pair", zap.Error(err))
		ErrorJSON(w, "failed to generate token pair")
		return
	}

	// 5. return success response
	response := CommonResponse{
		Status:  http.StatusOK,
		Message: "Successfully refreshed token pair",
		Data:    tokenPair,
	}

	_ = WriteJSON(w, http.StatusOK, response)
}

// CreateTBA godoc
//
//	@Tags			users
//	@Summary		creates DID bounded TBA for user
//	@Description	creates DID bounded TBA for user

// @Router			/users/create-tba [post]
func (c *UserCtrl) createTBA(w http.ResponseWriter, r *http.Request) {
	// 1. get user from context
	jwtUser, err := c.jwt.FromContext(r.Context())
	if err != nil {
		zap.L().Error("failed to get user from context", zap.Error(err))
		ErrorJSON(w, "user not found")
		return
	}

	fn := func(ctx context.Context) (interface{}, error) {
		// 2. check if user already has TBA
		u, err := c.user.FindUserByDID(ctx, jwtUser.DID)
		if err != nil {
			return nil, err
		}

		if u.TBAAddress != "" {
			return nil, user.ErrTBAAlreadyExists
		}

		// TODO

		// 3. create TBA

		// 4. update user with TBA

		// 5. return success response

		return nil, nil
	}

	_, err = c.tx.Exec(r.Context(), fn)
	if err != nil {
		if err == user.ErrTBAAlreadyExists {
			ErrorJSON(w, "TBA already exists")
			return
		}

		zap.L().Error("failed to create TBA", zap.Error(err))
		ErrorJSON(w, "failed to create TBA")
		return
	}

	response := CommonResponse{
		Status:  http.StatusCreated,
		Message: "Successfully created TBA",
	}

	_ = WriteJSON(w, http.StatusCreated, response)
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
func (c *UserCtrl) profile(w http.ResponseWriter, r *http.Request) {
	// 1. get user from context
	jwtUser, err := c.jwt.FromContext(r.Context())
	if err != nil {
		zap.L().Error("failed to get user from context", zap.Error(err))
		ErrorJSON(w, "user not found")
		return
	}

	// 2. get user from db
	u, err := c.user.FindUserByDID(r.Context(), jwtUser.DID)
	if err != nil {
		zap.L().Error("failed to find user", zap.Error(err))
		ErrorJSON(w, "failed to find user")
		return
	}

	resp := CommonResponse{
		Status:  http.StatusOK,
		Message: "Successfully fetched user profile",
		Data:    u,
	}

	// 3. return user as json response
	_ = WriteJSON(w, http.StatusOK, resp)
}

// PurchasedTickets godoc
//
//	@Tags			users
//	@Summary		returns purchased tickets
//	@Description	returns purchased tickets
//	@Produce		json

// @Router			/users/purchased-tickets [get]
func (c *UserCtrl) purchasedTickets(w http.ResponseWriter, r *http.Request) {
	// 1. get user from context
	_, err := c.jwt.FromContext(r.Context())
	if err != nil {
		zap.L().Error("failed to get user from context", zap.Error(err))
		ErrorJSON(w, "user not found")
		return
	}

	// 2. get purchased tickets from db

	// 3. return purchased tickets as json response
}

// IssuedTickets godoc
//
//	@Tags			users
//	@Summary		returns issued tickets
//	@Description	returns issued tickets
//	@Produce		json

// @Router			/users/issued-tickets [get]
func (c *UserCtrl) issuedTickets(w http.ResponseWriter, r *http.Request) {
	// 1. get user from context
	_, err := c.jwt.FromContext(r.Context())
	if err != nil {
		zap.L().Error("failed to get user from context", zap.Error(err))
		ErrorJSON(w, "user not found")
		return
	}

	// 2. get issued tickets from db

	// 3. return issued tickets as json response
}
