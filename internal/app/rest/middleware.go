package rest

import (
	"net/http"
	"strings"

	"github.com/heroticket/internal/service/jwt"
	"go.uber.org/zap"
)

func TokenRequired(jwtSvc jwt.Service) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Vary", "Authorization")

			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				ErrorJSON(w, "authorization header not found", http.StatusUnauthorized)
				return
			}

			headerParts := strings.Split(authHeader, " ")
			if len(headerParts) != 2 {
				ErrorJSON(w, "invalid authorization header", http.StatusUnauthorized)
				return
			}

			if headerParts[0] != "Bearer" {
				ErrorJSON(w, "invalid authorization header", http.StatusUnauthorized)
				return
			}

			token := headerParts[1]

			jwtUser, err := jwtSvc.VerifyToken(token, jwt.TokenRoleAccess)
			if err != nil {
				zap.L().Error("failed to validate access token", zap.Error(err))
				ErrorJSON(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			r = r.WithContext(jwtSvc.NewContext(r.Context(), *jwtUser))

			next.ServeHTTP(w, r)
		})
	}
}

func TokenCheck(jwtSvc jwt.Service) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Vary", "Authorization")

			check := func() {
				authHeader := r.Header.Get("Authorization")

				if authHeader == "" {
					return
				}

				headerParts := strings.Split(authHeader, " ")
				if len(headerParts) != 2 {
					return
				}

				if headerParts[0] != "Bearer" {
					return
				}

				token := headerParts[1]

				jwtUser, err := jwtSvc.VerifyToken(token, jwt.TokenRoleAccess)
				if err != nil {
					zap.L().Error("failed to validate access token", zap.Error(err))
					return
				}

				r = r.WithContext(jwtSvc.NewContext(r.Context(), *jwtUser))
			}

			check()

			next.ServeHTTP(w, r)
		})
	}
}
