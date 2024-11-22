package v1

import (
	"as4/config"
	"as4/internal/controller/http/v1/middlewares"
	"as4/internal/controller/http/v1/routes"
	auth_interface "as4/internal/usecase/interface"
	"as4/pkg/monitoring"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func NewRouter(l *zap.Logger, authUsecase auth_interface.Auth, config *config.Config) *mux.Router {
	r := mux.NewRouter()

	r.Use(middlewares.LoggerMiddleware(l))

	monitoring.RegisterMetrics()

	r.Use(middlewares.MetricsMiddleware)
	routes.NewMetricsRoute(r, config)
	v1 := r.PathPrefix("/v1").Subrouter()
	{
		routes.NewAuthRoute(v1, l, authUsecase, config)
		routes.NewHealthRoute(v1, l, config)
		routes.NewFormRoute(v1, l, config)
	}

	return r
}
