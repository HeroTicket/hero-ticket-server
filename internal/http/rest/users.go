package rest

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

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
	r.Post("/logout", c.logout)
	r.Post("/refresh-token", c.refreshToken)

	r.Group(func(r chi.Router) {
		r.Use(AccessTokenRequired(c.jwt))

		r.Post("/create-tba", nil)
		r.Get("/purchased-tickets", nil)
		r.Get("/issued-tickets", nil)
	})

	return r
}

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

	// 4. return login request as json response
	_ = WriteJSON(w, http.StatusOK, request)
}

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
	authResponse, err := c.did.LoginCallback(r.Context(), sessionId, string(tokenBytes))
	if err != nil {
		zap.L().Error("failed to handle login callback", zap.Error(err))
		ErrorJSON(w, "failed to handle login callback", http.StatusInternalServerError)
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
		DID:     u.DID,
		Address: u.WalletAddress,
	}

	tokenPair, err := c.jwt.GenerateTokenPair(jwtUser)
	if err != nil {
		zap.L().Error("failed to generate token pair", zap.Error(err))
		ErrorJSON(w, "failed to generate token pair", http.StatusInternalServerError)
		return
	}

	// 7. set token pair as cookie
	accessCookie := c.newCookie("access_token", tokenPair.AccessToken, "localhost", tokenPair.AccessTokenExpiry)
	refreshCookie := c.newCookie("refresh_token", tokenPair.RefreshToken, "localhost", tokenPair.RefreshTokenExpiry)

	http.SetCookie(w, accessCookie)
	http.SetCookie(w, refreshCookie)

	// 8. return success response
	response := CommonResponse{
		Status:  http.StatusOK,
		Message: fmt.Sprintf("User with ID %s Successfully authenticated", userDID),
	}

	_ = WriteJSON(w, http.StatusOK, response)
}

func (c *UserCtrl) logout(w http.ResponseWriter, r *http.Request) {
	// 1. generate expired cookies
	accessCookie := c.newCookie("access_token", "", "", -time.Hour)
	refreshCookie := c.newCookie("refresh_token", "", "", -time.Hour)

	http.SetCookie(w, accessCookie)
	http.SetCookie(w, refreshCookie)

	// 2. return success response
	response := CommonResponse{
		Status:  http.StatusOK,
		Message: "Successfully logged out",
	}

	_ = WriteJSON(w, http.StatusOK, response)
}

func (c *UserCtrl) refreshToken(w http.ResponseWriter, r *http.Request) {
	// 1. get refresh token from cookie
	refreshCookie, err := r.Cookie("refresh_token")
	if err != nil {
		zap.L().Error("refresh token not found in cookie", zap.Error(err))
		ErrorJSON(w, "refresh token not found in cookie")
		return
	}

	// 2. validate refresh token
	refreshToken := refreshCookie.Value

	// 3. validate refresh token
	jwtUser, err := c.jwt.VerifyToken(refreshToken, jwt.TokenRoleRefresh)
	if err != nil {
		zap.L().Error("failed to validate refresh token", zap.Error(err))
		ErrorJSON(w, "failed to validate refresh token")
		return
	}

	// 4. get user from db
	u, err := c.user.FindUserByDID(r.Context(), jwtUser.DID)
	if err != nil {
		zap.L().Error("failed to find user", zap.Error(err))
		ErrorJSON(w, "failed to find user")
		return
	}

	newJwtUser := jwt.JWTUser{
		DID:     u.DID,
		Address: u.WalletAddress,
	}

	// 5. generate new token pair
	tokenPair, err := c.jwt.GenerateTokenPair(newJwtUser)
	if err != nil {
		zap.L().Error("failed to generate token pair", zap.Error(err))
		ErrorJSON(w, "failed to generate token pair")
		return
	}

	// 6. set token pair as cookie
	accessCookie := c.newCookie("access_token", tokenPair.AccessToken, "localhost", tokenPair.AccessTokenExpiry)
	refreshCookie = c.newCookie("refresh_token", tokenPair.RefreshToken, "localhost", tokenPair.RefreshTokenExpiry)

	http.SetCookie(w, accessCookie)
	http.SetCookie(w, refreshCookie)

	// 7. return success response
	response := CommonResponse{
		Status:  http.StatusOK,
		Message: "Successfully refreshed token pair",
	}

	_ = WriteJSON(w, http.StatusOK, response)
}

func (c *UserCtrl) newCookie(name, value, domain string, expiry time.Duration) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Domain:   domain,
		Path:     "/",
		Expires:  time.Now().Add(expiry),
		MaxAge:   int(expiry.Seconds()),
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		Secure:   true,
	}
}
