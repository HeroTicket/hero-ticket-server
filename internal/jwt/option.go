package jwt

import "time"

var (
	defaultIssuer      = "hero-ticket"
	defaultAudience    = "localhost"
	defaultTokenExpiry = time.Minute * 30
)

type Option func(*jwtService)

func WithDefaultOptions() Option {
	return func(s *jwtService) {
		s.issuer = defaultIssuer
		s.audience = defaultAudience
		s.tokenExpiry = defaultTokenExpiry
	}
}

func WithIssuer(issuer string) Option {
	return func(s *jwtService) {
		s.issuer = issuer
	}
}

func WithAudience(audience string) Option {
	return func(s *jwtService) {
		s.audience = audience
	}
}

func WithTokenExpiry(expiry time.Duration) Option {
	return func(s *jwtService) {
		s.tokenExpiry = expiry
	}
}
