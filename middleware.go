// pkg/server/middleware.go
package server

import (
	"log"
	"net/http"
	"time"
)

// Middleware represents a function that wraps an http.Handler
type Middleware func(http.Handler) http.Handler

// ResponseWriter extends http.ResponseWriter with additional functionality
type ResponseWriter struct {
	http.ResponseWriter
	Status int
	Size   int64
}

// Common middleware constructors
func NewLoggingMiddleware(logger *log.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			rw := &ResponseWriter{ResponseWriter: w}
			next.ServeHTTP(rw, r)

			logger.Printf(
				"%s %s %d %s",
				r.Method,
				r.RequestURI,
				rw.Status,
				time.Since(start),
			)
		})
	}
}

func NewRecoveryMiddleware(logger *log.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Printf("panic: %v", err)
					http.Error(w,
						http.StatusText(http.StatusInternalServerError),
						http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// NewAuthMiddleware sets up authentication
func NewAuthMiddleware(sessionManager SessionManager) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := sessionManager.GetSession(r)
			if err != nil || !session.Authenticated {
				http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// CombineMiddleware chains multiple middleware together
func CombineMiddleware(middleware []Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middleware) - 1; i >= 0; i-- {
			next = middleware[i](next)
		}
		return next
	}
}
