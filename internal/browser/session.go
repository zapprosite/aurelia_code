// Package browser provides browser automation using go-rod (CDP).
package browser

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

// Session represents a browser session with persistent state.
type Session struct {
	ID          string
	Browser     *Browser
	UserDataDir string
	CreatedAt   time.Time
	LastActive  time.Time
}

// SessionConfig holds session configuration.
type SessionConfig struct {
	// DataDir is the base directory for session data.
	DataDir string
	// MaxAge is the maximum age before session is invalidated.
	MaxAge time.Duration
	// IdleTimeout is the timeout for idle sessions.
	IdleTimeout time.Duration
}

// DefaultSessionConfig returns sensible defaults.
func DefaultSessionConfig() SessionConfig {
	return SessionConfig{
		DataDir:     "/tmp/aurelia-browser-sessions",
		MaxAge:     24 * time.Hour,
		IdleTimeout: 10 * time.Minute,
	}
}

// NewSession creates a new browser session.
func NewSession(cfg SessionConfig, browserCfg Config) (*Session, error) {
	// Create session ID
	sessionID := uuid.New().String()

	// Create user data dir
	userDataDir := filepath.Join(cfg.DataDir, sessionID)
	if err := os.MkdirAll(userDataDir, 0755); err != nil {
		return nil, fmt.Errorf("create user data dir: %w", err)
	}

	// Update browser config with user data dir
	browserCfg.UserDataDir = userDataDir

	// Create browser
	browser, err := New(browserCfg)
	if err != nil {
		os.RemoveAll(userDataDir)
		return nil, fmt.Errorf("create browser: %w", err)
	}

	return &Session{
		ID:          sessionID,
		Browser:     browser,
		UserDataDir: userDataDir,
		CreatedAt:   time.Now(),
		LastActive:  time.Now(),
	}, nil
}

// Touch updates the last active time.
func (s *Session) Touch() {
	s.LastActive = time.Now()
}

// IsExpired checks if the session has expired based on maxAge.
func (s *Session) IsExpired(maxAge time.Duration) bool {
	if maxAge > 0 && time.Since(s.CreatedAt) > maxAge {
		return true
	}
	return false
}

// IsIdle checks if the session is idle.
func (s *Session) IsIdle(timeout time.Duration) bool {
	return time.Since(s.LastActive) > timeout
}

// Close closes the session and cleans up resources.
func (s *Session) Close() error {
	if s.Browser != nil {
		s.Browser.Close()
	}

	// Clean up user data dir
	if s.UserDataDir != "" {
		os.RemoveAll(s.UserDataDir)
	}

	return nil
}

// SessionManager manages browser sessions.
type SessionManager struct {
	sessions   map[string]*Session
	config     SessionConfig
	browserCfg Config
}

// NewSessionManager creates a new session manager.
func NewSessionManager(browserCfg Config, config SessionConfig) *SessionManager {
	// Ensure data dir exists
	os.MkdirAll(config.DataDir, 0755)

	return &SessionManager{
		sessions:   make(map[string]*Session),
		config:     config,
		browserCfg: browserCfg,
	}
}

// Acquire gets or creates a session.
func (m *SessionManager) Acquire() (*Session, error) {
	// Find an existing non-expired session
	for _, session := range m.sessions {
		if !session.IsExpired(m.config.MaxAge) && !session.IsIdle(m.config.IdleTimeout) {
			session.Touch()
			return session, nil
		}
	}

	// Create new session
	session, err := NewSession(m.config, m.browserCfg)
	if err != nil {
		return nil, err
	}

	m.sessions[session.ID] = session
	return session, nil
}

// Release returns a session to the pool.
func (m *SessionManager) Release(session *Session) {
	if session != nil {
		session.Touch()
	}
}

// Close closes all sessions.
func (m *SessionManager) Close() error {
	for _, session := range m.sessions {
		session.Close()
	}
	m.sessions = make(map[string]*Session)
	return nil
}

// Cleanup removes expired sessions.
func (m *SessionManager) Cleanup() {
	for id, session := range m.sessions {
		if session.IsExpired(m.config.MaxAge) || session.IsIdle(m.config.IdleTimeout) {
			session.Close()
			delete(m.sessions, id)
		}
	}
}

// Count returns the number of active sessions.
func (m *SessionManager) Count() int {
	return len(m.sessions)
}
