package routes

import (
	"as4/config"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewMetricsRoute(router *mux.Router, config *config.Config) {
	router.Handle("/metrics", promhttp.Handler()).Methods("GET")
}
