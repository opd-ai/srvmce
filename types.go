// pkg/server/types.go
package server

import (
	"net/http"
)

// SessionManager handles user sessions
type SessionManager interface {
	GetSession(r *http.Request) (*Session, error)
	CreateSession(w http.ResponseWriter, userID string) error
	DestroySession(w http.ResponseWriter, r *http.Request) error
}

// Session represents a user session
type Session struct {
	ID            string
	UserID        string
	Authenticated bool
	CreatedAt     int64
}

// DefaultSessionManager implements SessionManager
type DefaultSessionManager struct {
	// In a production environment, you'd want to use a proper session store
	sessions map[string]*Session
}

func NewSessionManager() SessionManager {
	return &DefaultSessionManager{
		sessions: make(map[string]*Session),
	}
}
