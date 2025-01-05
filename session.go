// pkg/server/session.go
package server

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"
)

// Implementation of DefaultSessionManager
func (sm *DefaultSessionManager) GetSession(r *http.Request) (*Session, error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return nil, err
	}

	session, exists := sm.sessions[cookie.Value]
	if !exists {
		return nil, fmt.Errorf("session not found")
	}

	return session, nil
}

func (sm *DefaultSessionManager) CreateSession(w http.ResponseWriter, userID string) error {
	sessionID, err := generateSessionID()
	if err != nil {
		return err
	}

	session := &Session{
		ID:            sessionID,
		UserID:        userID,
		Authenticated: true,
		CreatedAt:     time.Now().Unix(),
	}

	sm.sessions[sessionID] = session

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	return nil
}

func (sm *DefaultSessionManager) DestroySession(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie("session")
	if err != nil {
		return nil
	}

	delete(sm.sessions, cookie.Value)

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	return nil
}

func generateSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
