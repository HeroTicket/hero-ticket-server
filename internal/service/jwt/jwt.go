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
	AccessToken        string        `json:"accessToken"`
	RefreshToken       string        `json:"refreshToken"`
	AccessTokenExpiry  time.Duration `json:"accessTokenExpiry"`
	RefreshTokenExpiry time.Duration `json:"refreshTokenExpiry"`
}

type JWTUser struct {
	ID string `json:"id"`
}

type JWTUserKey struct{}
