package middlewares

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func LoggerMiddleware(l *zap.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l.Info("Request", zap.String("method", r.Method), zap.String("url", r.URL.Path))
			next.ServeHTTP(w, r)
		})
	}
}
