package middlewares

import "net/http"

// SecurityHeadersMiddleware adds common security headers to all responses.
func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Prevent Clickjacking by disallowing iframe embedding
		w.Header().Set("X-Frame-Options", "DENY")

		// Mitigate XSS attacks by controlling where resources can be loaded from
		w.Header().Set("Content-Security-Policy", "default-src 'self'")

		// Prevent MIME sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// Strict HTTPS enforcement (browsers will always use HTTPS)
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// Enable XSS protection in older browsers
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		next.ServeHTTP(w, r)
	})
}
