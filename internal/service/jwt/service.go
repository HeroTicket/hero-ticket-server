package jwt

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Service interface {
	GenerateToken(u JWTUser) (*JwtToken, error)
	VerifyToken(rawToken string) (*JWTUser, error)
	NewContext(ctx context.Context, u JWTUser) context.Context
	FromContext(ctx context.Context) (*JWTUser, error)
}

type jwtService struct {
	issuer      string
	audience    string
	secretKey   string
	tokenExpiry time.Duration
}

func New(secretKey string, opts ...Option) Service {
	svc := &jwtService{
		secretKey: secretKey,
	}

	WithDefaultOptions()(svc)

	for _, opt := range opts {
		opt(svc)
	}

	return svc
}

func (s *jwtService) GenerateToken(u JWTUser) (*JwtToken, error) {
	token, err := s.generateToken(s.getTokenClaims(u), s.secretKey)
	if err != nil {
		return nil, err
	}

	return &JwtToken{
		Token:       token,
		TokenExpiry: s.tokenExpiry,
	}, nil
}

func (s *jwtService) VerifyToken(rawToken string) (*JWTUser, error) {
	// get key function based on token role
	fn := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}

		return []byte(s.secretKey), nil
	}

	token, err := jwt.Parse(rawToken, fn)
	if err != nil {
		return nil, err
	}

	var u JWTUser

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
			u.Identifier = sub
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

func (s *jwtService) getTokenClaims(u JWTUser) jwt.MapClaims {
	return jwt.MapClaims{
		"iss": s.issuer,
		"aud": s.audience,
		"sub": u.Identifier,
		"exp": time.Now().Add(s.tokenExpiry).UTC().Unix(),
		"iat": time.Now().UTC().Unix(),
		"typ": "JWT",
	}
}
