package ssh

import (
	"sync"
	"time"
)

// SessionManager manages active sessions
type SessionManager struct {
	sessions map[string]*Session
	mu       sync.RWMutex
}

// NewSessionManager creates a new session manager
func NewSessionManager() *SessionManager {
	sm := &SessionManager{
		sessions: make(map[string]*Session),
	}
	// Start cleanup loop
	go sm.cleanupLoop()
	return sm
}

// Add adds a session
func (sm *SessionManager) Add(s *Session) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.sessions[s.ID] = s
}

// Get returns a session by ID
func (sm *SessionManager) Get(id string) *Session {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.sessions[id]
}

// Remove removes a session
func (sm *SessionManager) Remove(id string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if s, ok := sm.sessions[id]; ok {
		s.Close()
		delete(sm.sessions, id)
	}
}

// cleanupLoop periodically removes inactive sessions
func (sm *SessionManager) cleanupLoop() {
	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		sm.cleanup()
	}
}

func (sm *SessionManager) cleanup() {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	
	now := time.Now()
	for id, s := range sm.sessions {
		// Check if active or detached
		s.WsLock.Lock()
		active := s.WsConn != nil
		s.WsLock.Unlock()

		if !active {
			// Detached mode
			// If KeepAliveUntil is set and specific time holds
			if !s.KeepAliveUntil.IsZero() {
				if now.After(s.KeepAliveUntil) {
					s.Close()
					delete(sm.sessions, id)
				}
			} else {
				// No keepalive set, remove immediately (or treat as abandon)
				// Actually handler should have removed it if MaxKeepAlive was 0.
				// But as a fallback, we remove here.
				s.Close()
				delete(sm.sessions, id)
			}
		} else {
			// Active mode
			// Close session if inactive for more than 1 hour (auto-logout idle users)
			if now.Sub(s.LastActive) > time.Hour {
				s.Close()
				delete(sm.sessions, id)
			}
		}
	}
}
