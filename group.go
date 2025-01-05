// pkg/server/group.go
package server

import (
	"net/http"
	"path"
)

// Group represents a route group with shared middleware and prefix
type Group struct {
	server     *HTTPServer
	prefix     string
	middleware []Middleware
}

func NewGroup(server *HTTPServer, prefix string, middleware ...Middleware) *Group {
	return &Group{
		server:     server,
		prefix:     prefix,
		middleware: middleware,
	}
}

func (g *Group) Handle(pattern string, handler http.Handler, middleware ...Middleware) {
	fullPath := path.Join(g.prefix, pattern)
	allMiddleware := append(g.middleware, middleware...)
	g.server.Handle(fullPath, handler, allMiddleware...)
}

func (g *Group) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request),
	middleware ...Middleware,
) {
	g.Handle(pattern, http.HandlerFunc(handler), middleware...)
}

func (g *Group) Group(prefix string, middleware ...Middleware) RouterGroup {
	fullPrefix := path.Join(g.prefix, prefix)
	allMiddleware := append(g.middleware, middleware...)
	return NewGroup(g.server, fullPrefix, allMiddleware...)
}
