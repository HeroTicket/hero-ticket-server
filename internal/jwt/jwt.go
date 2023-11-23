package jwt

import (
	"errors"
	"time"
)

var (
	ErrInvalidContext          = errors.New("invalid context")
	ErrInvalidToken            = errors.New("invalid token")
	ErrInvalidTokenRole        = errors.New("invalid token role")
	ErrInvalidSigningMethod    = errors.New("invalid token signing method")
	ErrAccessTokenKeyRequired  = errors.New("access token key is required")
	ErrRefreshTokenKeyRequired = errors.New("refresh token key is required")
)

type TokenRole uint8

const (
	TokenRoleAccess TokenRole = iota + 1
	TokenRoleRefresh
)

type TokenPair struct {
	AccessToken        string        `json:"access_token"`
	RefreshToken       string        `json:"refresh_token"`
	AccessTokenExpiry  time.Duration `json:"access_token_expiry"`
	RefreshTokenExpiry time.Duration `json:"refresh_token_expiry"`
}

type JWTUser struct {
	DID string `json:"did"`
}

type JWTUserKey struct{}
