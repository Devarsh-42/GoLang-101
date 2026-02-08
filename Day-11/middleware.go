package main

import (
	"log/slog"
	"net/http"
	"time"
)

// ============================================================
// middleware.go - HTTP middleware functions
// Middleware wraps handlers to add extra behavior (logging, etc.)
// ============================================================

// loggingMiddleware logs every incoming HTTP request with method, path, and duration
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now() // Record when the request started

		// Call the actual handler
		next.ServeHTTP(w, r)

		// Log the request details after it completes
		slog.Info("HTTP Request",
			"method", r.Method,
			"path", r.URL.Path,
			"duration", time.Since(start).String(),
		)
	})
}

// jsonContentTypeMiddleware sets Content-Type to application/json for API routes
func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json") // Tell the client it's JSON
		next.ServeHTTP(w, r)                               // Call the actual handler
	})
}
