package jwt

import "time"

var (
	defaultIssuer             = "hero-ticket"
	defaultAudience           = "localhost"
	defaultAccessTokenExpiry  = time.Minute * 15
	defaultRefreshTokenExpiry = time.Hour * 24
)

type Option func(*Service)

func WithDefaultOptions() Option {
	return func(s *Service) {
		s.issuer = defaultIssuer
		s.audience = defaultAudience
		s.accessTokenExpiry = defaultAccessTokenExpiry
		s.refreshTokenExpiry = defaultRefreshTokenExpiry
	}
}

func WithIssuer(issuer string) Option {
	return func(s *Service) {
		s.issuer = issuer
	}
}

func WithAudience(audience string) Option {
	return func(s *Service) {
		s.audience = audience
	}
}

func WithAccessTokenExpiry(expiry time.Duration) Option {
	return func(s *Service) {
		s.accessTokenExpiry = expiry
	}
}

func WithRefreshTokenExpiry(expiry time.Duration) Option {
	return func(s *Service) {
		s.refreshTokenExpiry = expiry
	}
}
