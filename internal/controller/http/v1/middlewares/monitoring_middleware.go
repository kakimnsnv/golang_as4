package middlewares

import (
	"as4/pkg/monitoring"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		rec := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rec, r)

		path := r.URL.Path

		monitoring.RequestCounter.With(prometheus.Labels{
			"path":   path,
			"method": r.Method,
			"status": strconv.Itoa(rec.statusCode),
		}).Inc()

		if rec.statusCode >= http.StatusBadRequest {
			monitoring.ErrorCounter.With(prometheus.Labels{
				"path":   path,
				"status": strconv.Itoa(rec.statusCode),
			}).Inc()
		}

		monitoring.RequestDuration.With(prometheus.Labels{
			"path": path,
		}).Inc()
	})
}

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}
