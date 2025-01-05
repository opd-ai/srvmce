package server

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"path"
	"sync"
)

//go:embed assets/templates/* assets/static/admin/*
var assets embed.FS

// Initialize templates
func (s *HTTPServer) initTemplates() error {
	templates, err := template.ParseFS(assets,
		"assets/templates/admin/*.html",
		"assets/templates/errors/*.html")
	if err != nil {
		return fmt.Errorf("failed to parse templates: %w", err)
	}
	s.templates = templates
	return nil
}

// Serve static files
func (s *HTTPServer) serveStaticFiles() {
	fs := http.FileServer(http.FS(assets))
	s.Handle("/admin/static/", http.StripPrefix("/admin/static/",
		fs))
}

// Server represents the main web server interface
type Server interface {
	Serve(l net.Listener) error
	Shutdown(ctx context.Context) error
	Handle(pattern string, handler http.Handler, middleware ...Middleware)
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request), middleware ...Middleware)
	Group(prefix string, middleware ...Middleware) RouterGroup
}

// Router handles HTTP request routing
type Router struct {
	mu       sync.RWMutex
	Routes   map[string]http.Handler
	NotFound http.Handler
}

// RouterGroup represents a group of routes with shared middleware
type RouterGroup interface {
	Handle(pattern string, handler http.Handler, middleware ...Middleware)
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request), middleware ...Middleware)
	Group(prefix string, middleware ...Middleware) RouterGroup
}

// Config holds the server configuration
type Config struct {
	RootDir        string
	AdminPath      string
	TinyMCEKey     string
	AdminUser      string
	AdminPassword  string
	AdminTemplates *template.Template
}

func NewServer(cfg Config) (Server, error) {
	if err := validateConfig(cfg); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	s := &HTTPServer{
		Router: NewRouter(),
		Config: cfg,
	}

	if err := s.setupAdminRoutes(); err != nil {
		return nil, fmt.Errorf("failed to setup admin routes: %w", err)
	}

	return s, nil
}

func NewRouter() *Router {
	return &Router{
		Routes:   make(map[string]http.Handler),
		NotFound: http.NotFoundHandler(),
	}
}

// Implementation methods for HTTPServer
func (s *HTTPServer) Serve(l net.Listener) error {
	s.httpServer = &http.Server{
		Handler: s.Router,
	}
	return s.httpServer.Serve(l)
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	if s.httpServer == nil {
		return nil
	}
	return s.httpServer.Shutdown(ctx)
}

func (s *HTTPServer) Handle(pattern string, handler http.Handler, middleware ...Middleware) {
	s.mu.Lock()
	defer s.mu.Unlock()

	chain := CombineMiddleware(append(s.Middlewares, middleware...))
	s.Router.Handle(pattern, chain(handler))
}

// Router implementation
func (r *Router) Handle(pattern string, handler http.Handler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Routes[path.Clean(pattern)] = handler
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mu.RLock()
	handler, exists := r.Routes[path.Clean(req.URL.Path)]
	r.mu.RUnlock()

	if !exists {
		r.NotFound.ServeHTTP(w, req)
		return
	}
	handler.ServeHTTP(w, req)
}
