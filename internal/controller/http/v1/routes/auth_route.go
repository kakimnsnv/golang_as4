package routes

import (
	"as4/config"
	"as4/internal/entity"
	auth_interface "as4/internal/usecase/interface"
	"as4/pkg/validator"
	"encoding/json"
	"fmt"
	"net/http"

	v "github.com/go-playground/validator/v10"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type authRoute struct {
	auth_usecase auth_interface.Auth
	logger       *zap.Logger
	config       *config.Config
}

func NewAuthRoute(
	router *mux.Router,
	logger *zap.Logger,
	auth_usecase auth_interface.Auth,
	config *config.Config,
) {
	r := authRoute{
		auth_usecase: auth_usecase,
		logger:       logger,
		config:       config,
	}

	router.HandleFunc("/csrf-token", r.getCSRFToken).Methods("GET")

	s := router.PathPrefix("/auth").Subrouter()
	{
		s.HandleFunc("/login", r.login).Methods("POST")
		s.HandleFunc("/register", r.register).Methods("POST")
	}
}

func (a *authRoute) login(w http.ResponseWriter, r *http.Request) {
	var loginRequest entity.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		a.logger.Info("failed to decode request", zap.Error(err))
		http.Error(w, "failed to decode request", http.StatusBadRequest)
		return
	}

	if err := validator.ValidateStruct(loginRequest); err != nil {
		if validationErrors, ok := err.(v.ValidationErrors); ok {
			customErrors := validator.FormatValidationError(validationErrors)
			a.logger.Info("failed to validate request", zap.Error(err))
			http.Error(w, fmt.Sprintf("Validation errors: %v", customErrors), http.StatusBadRequest)
			return
		}
		a.logger.Warn("validator returned not type of ValidationErrors", zap.Error(err))
	}

	authResponse, err := a.auth_usecase.Login(r.Context(), loginRequest.Email, loginRequest.Password, a.config)
	if err != nil {
		a.logger.Info("failed to login", zap.Error(err))
		http.Error(w, "failed to login", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(authResponse)
}

func (a *authRoute) register(w http.ResponseWriter, r *http.Request) {
	var registerRequest entity.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&registerRequest); err != nil {
		a.logger.Info("failed to decode request", zap.Error(err))
		http.Error(w, "failed to decode request", http.StatusBadRequest)
		return
	}

	if err := validator.ValidateStruct(registerRequest); err != nil {
		if validationErrors, ok := err.(v.ValidationErrors); ok {
			customErrors := validator.FormatValidationError(validationErrors)
			a.logger.Info("failed to validate request", zap.Error(err))
			http.Error(w, fmt.Sprintf("Validation errors: %v", customErrors), http.StatusBadRequest)
			return
		}
		a.logger.Warn("validator returned not type of ValidationErrors", zap.Error(err))
	}

	authResponse, err := a.auth_usecase.Register(r.Context(), registerRequest.Email, registerRequest.Password, registerRequest.Username, a.config)
	if err != nil {
		a.logger.Info("failed to register", zap.Error(err))
		http.Error(w, "failed to register", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(authResponse)
}

func (a *authRoute) getCSRFToken(w http.ResponseWriter, r *http.Request) {
	token := csrf.Token(r)
	w.Header().Set("X-Csrf-Token", token)
	w.WriteHeader(http.StatusOK)
}
