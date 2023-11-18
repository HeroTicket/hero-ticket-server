package rest

import (
	"net/http"

	"github.com/heroticket/internal/jwt"
)

func AccessTokenRequired(jwtSvc jwt.Service) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("access_token")
			if err != nil {
				ErrorJSON(w, "access token not found in cookie", http.StatusUnauthorized)
				return
			}

			accessToken := cookie.Value

			jwtUser, err := jwtSvc.VerifyToken(accessToken, jwt.TokenRoleAccess)
			if err != nil {
				ErrorJSON(w, "failed to validate access token", http.StatusUnauthorized)
				return
			}

			r = r.WithContext(jwtSvc.NewContext(r.Context(), *jwtUser))

			next.ServeHTTP(w, r)
		})
	}
}
