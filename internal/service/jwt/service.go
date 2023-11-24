package jwt

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Service interface {
	GenerateTokenPair(u JWTUser) (*TokenPair, error)
	VerifyToken(token string, role TokenRole) (*JWTUser, error)
	NewContext(ctx context.Context, u JWTUser) context.Context
	FromContext(ctx context.Context) (*JWTUser, error)
}

type jwtService struct {
	issuer             string
	audience           string
	accessTokenKey     string
	refreshTokenKey    string
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
}

func New(accessTokenKey, refreshTokenKey string, opts ...Option) Service {
	svc := &jwtService{
		accessTokenKey:  accessTokenKey,
		refreshTokenKey: refreshTokenKey,
	}

	WithDefaultOptions()(svc)

	for _, opt := range opts {
		opt(svc)
	}

	return svc
}

// GenerateTokenPair generates a pair of access and refresh tokens
func (s *jwtService) GenerateTokenPair(u JWTUser) (*TokenPair, error) {
	// generate access token
	accessToken, err := s.generateToken(s.getAccessTokenClaims(u), s.accessTokenKey)
	if err != nil {
		return nil, err
	}

	// generate refresh token
	refreshToken, err := s.generateToken(s.getRefreshTokenClaims(u), s.refreshTokenKey)
	if err != nil {
		return nil, err
	}

	// return token pair with configured expiry
	return &TokenPair{
		AccessToken:        accessToken,
		RefreshToken:       refreshToken,
		AccessTokenExpiry:  s.accessTokenExpiry,
		RefreshTokenExpiry: s.refreshTokenExpiry,
	}, nil
}

// VerifyToken verifies a token and returns the user info if the token is valid
func (s *jwtService) VerifyToken(tokenString string, role TokenRole) (*JWTUser, error) {
	// get key function based on token role
	fn, err := s.getKeyfunc(role)
	if err != nil {
		return nil, err
	}

	// parse token
	token, err := jwt.Parse(tokenString, fn)
	if err != nil {
		return nil, err
	}

	var u JWTUser

	// validate token
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// check issuer and audience
		if iss, ok := claims["iss"].(string); !ok || iss != s.issuer {
			return nil, ErrInvalidToken
		}

		if aud, ok := claims["aud"].(string); !ok || aud != s.audience {
			return nil, ErrInvalidToken
		}

		// check subject
		if sub, ok := claims["sub"].(string); ok {
			u.ID = sub
		}

		// return jwt user
		return &u, nil
	}

	// invalid token
	return nil, ErrInvalidToken
}

// NewContext returns a new context with the user info
func (s *jwtService) NewContext(ctx context.Context, u JWTUser) context.Context {
	return context.WithValue(ctx, JWTUserKey{}, u)
}

// FromContext returns the user info from the context
func (s *jwtService) FromContext(ctx context.Context) (*JWTUser, error) {
	u, ok := ctx.Value(JWTUserKey{}).(JWTUser)
	if !ok {
		return nil, ErrInvalidContext
	}

	return &u, nil
}

// generateToken generates a token
func (s *jwtService) generateToken(claims jwt.MapClaims, key string) (string, error) {
	// create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// sign token
	signedToken, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	// return signed token
	return signedToken, nil
}

// getKeyfunc returns a key function based on token role
func (s *jwtService) getKeyfunc(role TokenRole) (jwt.Keyfunc, error) {
	var key string

	// get key based on token role
	switch role {
	case TokenRoleAccess:
		key = s.accessTokenKey
	case TokenRoleRefresh:
		key = s.refreshTokenKey
	default:
		// invalid token role
		return nil, ErrInvalidTokenRole
	}

	// return key function
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}

		return []byte(key), nil
	}, nil
}

// getAccessTokenClaims returns the claims for access token
func (s *jwtService) getAccessTokenClaims(u JWTUser) jwt.MapClaims {
	// access token contains email claim
	return jwt.MapClaims{
		"iss": s.issuer,
		"aud": s.audience,
		"sub": u.ID,
		"exp": time.Now().Add(s.accessTokenExpiry).UTC().Unix(),
		"iat": time.Now().UTC().Unix(),
		"typ": "JWT",
	}
}

// getRefreshTokenClaims returns the claims for refresh token
func (s *jwtService) getRefreshTokenClaims(u JWTUser) jwt.MapClaims {
	// refresh token does not contain email claim
	return jwt.MapClaims{
		"iss": s.issuer,
		"aud": s.audience,
		"sub": u.ID,
		"exp": time.Now().Add(s.refreshTokenExpiry).UTC().Unix(),
		"iat": time.Now().UTC().Unix(),
	}
}
