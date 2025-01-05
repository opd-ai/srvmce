// pkg/server/server.go
package server

import (
	"fmt"
	"net/http"
	"sync"

	html "html/template"
)

type HTTPServer struct {
	Router      *Router
	Config      Config
	httpServer  *http.Server
	mu          sync.RWMutex
	Middlewares []Middleware
	templates   *html.Template
}

// validateConfig checks if the server configuration is valid
func validateConfig(cfg Config) error {
	if cfg.RootDir == "" {
		return fmt.Errorf("root directory is required")
	}
	if cfg.AdminUser == "" || cfg.AdminPassword == "" {
		return fmt.Errorf("admin credentials are required")
	}
	return nil
}

// setupAdminRoutes initializes all admin-related routes
func (s *HTTPServer) setupAdminRoutes() error {
	adminGroup := s.Group(s.Config.AdminPath)

	// Setup admin routes
	adminGroup.HandleFunc("/login", s.handleAdminLogin)
	adminGroup.HandleFunc("/logout", s.handleAdminLogout)

	// Protected admin routes
	protected := adminGroup.Group("", NewAuthMiddleware(NewSessionManager()))
	protected.HandleFunc("/dashboard", s.handleAdminDashboard)
	protected.HandleFunc("/files", s.handleFileManager)

	return nil
}

// Admin handler methods
func (s *HTTPServer) handleAdminLogin(w http.ResponseWriter, r *http.Request) {
	// Implementation details
}

func (s *HTTPServer) handleAdminLogout(w http.ResponseWriter, r *http.Request) {
	// Implementation details
}

func (s *HTTPServer) handleAdminDashboard(w http.ResponseWriter, r *http.Request) {
	// Implementation details
}

func (s *HTTPServer) handleFileManager(w http.ResponseWriter, r *http.Request) {
	// Implementation details
}

// Group implements the Server interface
func (s *HTTPServer) Group(prefix string, middleware ...Middleware) RouterGroup {
	return NewGroup(s, prefix, middleware...)
}

func (s *HTTPServer) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request),
	middleware ...Middleware,
) {
	s.Handle(pattern, http.HandlerFunc(handler), middleware...)
}
