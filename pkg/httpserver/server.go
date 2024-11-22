package httpserver

import (
	"as4/config"
	"as4/internal/controller/http/v1/middlewares"
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewHTTPSServer(lc fx.Lifecycle, handler http.Handler, config *config.Config, l *zap.Logger) *http.Server {

	secureHandler := middlewares.SecurityHeadersMiddleware(handler)

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	srv := &http.Server{
		Addr:      config.Port,
		Handler:   secureHandler,
		TLSConfig: tlsConfig,
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := srv.ListenAndServeTLS("server.crt", "server.key"); err != nil && err != http.ErrServerClosed {
					l.Fatal("Failer to start HTTPS server", zap.Error(err))
				}
			}()
			l.Info("Starting HTTPS server")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()
			l.Info("Stopping HTTPS server")
			return srv.Shutdown(shutdownCtx)
		},
	})
	return srv
}
