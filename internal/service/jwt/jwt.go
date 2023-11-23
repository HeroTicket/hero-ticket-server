package jwt

import (
	"errors"
	"time"
)

var (
	ErrInvalidContext       = errors.New("invalid context")
	ErrInvalidToken         = errors.New("invalid token")
	ErrInvalidSigningMethod = errors.New("invalid token signing method")
)

type JwtToken struct {
	Token       string        `json:"token"`
	TokenExpiry time.Duration `json:"tokenExpiry"`
}

type JWTUser struct {
	Identifier string `json:"identifier"`
}

type JWTUserKey struct{}
