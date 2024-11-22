package middlewares

import (
	"as4/config"
	"as4/pkg/auth"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func AuthMiddleware(config *config.Config) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr := r.Header.Get("Authorization")
			if tokenStr == "" {
				http.Error(w, "no token provided", http.StatusBadRequest)
				return
			}
			token, err := auth.ValidateJWT(tokenStr, config)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			if token.ExpiresAt == nil {
				http.Error(w, "invalid token: no expiration date provided", http.StatusUnauthorized)
				return
			}

			if time.Now().Unix() > token.ExpiresAt.Unix() {
				http.Error(w, "token expired", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
