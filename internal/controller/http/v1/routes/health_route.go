package routes

import (
	"as4/config"
	"as4/internal/controller/http/v1/middlewares"
	"as4/pkg/auth"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type healthRoute struct {
	l      *zap.Logger
	config *config.Config
}

func NewHealthRoute(router *mux.Router, l *zap.Logger, config *config.Config) {
	r := &healthRoute{l: l, config: config}

	s := router.PathPrefix("/health").Subrouter()
	{
		s.Use(middlewares.AuthMiddleware(config))
		s.HandleFunc("/", r.Health).Methods("GET")
	}
}

func (h *healthRoute) Health(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		http.Error(w, "no token provided", http.StatusBadRequest)
		return
	}
	token, err := auth.ValidateJWT(tokenStr, h.config)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}
	message := "OK"
	if token.Role == "admin" {
		message = "OK admin"
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
}
