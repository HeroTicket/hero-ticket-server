package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	issuer             string
	audience           string
	accessTokenKey     string
	refreshTokenKey    string
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
}

func New(accessTokenKey, refreshTokenKey string, opts ...Option) (*Service, error) {
	svc := &Service{
		accessTokenKey:  accessTokenKey,
		refreshTokenKey: refreshTokenKey,
	}

	WithDefaultOptions()(svc)

	for _, opt := range opts {
		opt(svc)
	}

	return svc.validate()
}

// GenerateTokenPair generates a pair of access and refresh tokens
func (s *Service) GenerateTokenPair(u JWTUser) (*TokenPair, error) {
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
func (s *Service) VerifyToken(tokenString string, role TokenRole) (*JWTUser, error) {
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
			u.DID = sub
		}

		// only access token has address claim
		if address, ok := claims["address"].(string); ok {
			u.Address = address
		}

		// return jwt user
		return &u, nil
	}

	// invalid token
	return nil, ErrInvalidToken
}

// generateToken generates a token
func (s *Service) generateToken(claims jwt.MapClaims, key string) (string, error) {
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
func (s *Service) getKeyfunc(role TokenRole) (jwt.Keyfunc, error) {
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
func (s *Service) getAccessTokenClaims(u JWTUser) jwt.MapClaims {
	// access token contains email claim
	return jwt.MapClaims{
		"address": u.Address,
		"iss":     s.issuer,
		"aud":     s.audience,
		"sub":     u.DID,
		"exp":     time.Now().Add(s.accessTokenExpiry).UTC().Unix(),
		"iat":     time.Now().UTC().Unix(),
		"typ":     "JWT",
	}
}

// getRefreshTokenClaims returns the claims for refresh token
func (s *Service) getRefreshTokenClaims(u JWTUser) jwt.MapClaims {
	// refresh token does not contain email claim
	return jwt.MapClaims{
		"iss": s.issuer,
		"aud": s.audience,
		"sub": u.DID,
		"exp": time.Now().Add(s.refreshTokenExpiry).UTC().Unix(),
		"iat": time.Now().UTC().Unix(),
	}
}

// validate validates the service configuration
func (s *Service) validate() (*Service, error) {
	if s.accessTokenKey == "" {
		return nil, ErrAccessTokenKeyRequired
	}

	if s.refreshTokenKey == "" {
		return nil, ErrRefreshTokenKeyRequired
	}

	if s.issuer == "" {
		s.issuer = defaultIssuer
	}

	if s.audience == "" {
		s.audience = defaultAudience
	}

	if s.accessTokenExpiry == 0 {
		s.accessTokenExpiry = defaultAccessTokenExpiry
	}

	if s.refreshTokenExpiry == 0 {
		s.refreshTokenExpiry = defaultRefreshTokenExpiry
	}

	return s, nil
}
