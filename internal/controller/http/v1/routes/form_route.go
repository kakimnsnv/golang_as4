package routes

import (
	"as4/config"
	"as4/internal/entity"
	"as4/pkg/validator"
	"fmt"
	"net/http"

	v "github.com/go-playground/validator/v10"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type formRoute struct {
	l      *zap.Logger
	config *config.Config
}

func NewFormRoute(router *mux.Router, l *zap.Logger, config *config.Config) {
	r := &formRoute{l: l, config: config}

	router.HandleFunc("/form", r.formHandler).Methods("POST")

}

func (f *formRoute) formHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	formData := entity.FormData{
		Name:  r.FormValue("name"),
		Email: r.FormValue("email"),
	}

	if err := validator.ValidateStruct(formData); err != nil {
		if validationErrors, ok := err.(v.ValidationErrors); ok {
			customErrors := validator.FormatValidationError(validationErrors)
			f.l.Info("failed to validate request", zap.Error(err))
			http.Error(w, fmt.Sprintf("Validation errors: %v", customErrors), http.StatusBadRequest)
			return
		}
		f.l.Warn("validator returned not type of ValidationErrors", zap.Error(err))
	}

	fmt.Fprintf(w, "Name: %s\nEmail: %s\n", formData.Name, formData.Email)
}
