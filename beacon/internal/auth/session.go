package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	DefaultSessionDuration = 24 * time.Hour
	SessionIDBytes         = 32
	SessionCookieName      = "session_token"
	MaxSessionDuration     = 7 * 24 * time.Hour
)

type Session struct {
	ID         string
	UserID     string
	Scopes     Scopes
	CreatedAt  time.Time
	ExpiresAt  time.Time
	IPAddress  string
	UserAgent  string
	ParentID   string
	IsRotation bool
}

type SessionStore interface {
	Create(ctx context.Context, session Session) (string, error)
	Get(ctx context.Context, id string) (Session, error)
	Delete(ctx context.Context, id string) error
	Rotate(ctx context.Context, oldID string) (Session, string, error)
}

type InMemorySessionStore struct {
	mu       sync.RWMutex
	sessions map[string]Session
}

func NewInMemorySessionStore() *InMemorySessionStore {
	return &InMemorySessionStore{
		sessions: make(map[string]Session),
	}
}

func generateSessionID() (string, error) {
	b := make([]byte, SessionIDBytes)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate session ID: %w", err)
	}
	return hex.EncodeToString(b), nil
}

func generateRawSessionID() (string, string, error) {
	b := make([]byte, SessionIDBytes)
	if _, err := rand.Read(b); err != nil {
		return "", "", fmt.Errorf("failed to generate session ID: %w", err)
	}
	raw := hex.EncodeToString(b)
	hash := hashSessionID(raw)
	return raw, hash, nil
}

func hashSessionID(id string) string {
	h := sha256.Sum256([]byte(id))
	return hex.EncodeToString(h[:])
}

func (s *InMemorySessionStore) Create(ctx context.Context, session Session) (string, error) {
	rawID, hashedID, err := generateRawSessionID()
	if err != nil {
		return "", err
	}

	session.ID = hashedID
	if session.CreatedAt.IsZero() {
		session.CreatedAt = time.Now()
	}
	if session.ExpiresAt.IsZero() {
		session.ExpiresAt = time.Now().Add(DefaultSessionDuration)
	}
	if session.ExpiresAt.Sub(session.CreatedAt) > MaxSessionDuration {
		session.ExpiresAt = session.CreatedAt.Add(MaxSessionDuration)
	}

	s.mu.Lock()
	s.sessions[hashedID] = session
	s.mu.Unlock()

	return rawID, nil
}

func (s *InMemorySessionStore) Get(ctx context.Context, id string) (Session, error) {
	hashedID := hashSessionID(id)

	s.mu.RLock()
	session, ok := s.sessions[hashedID]
	s.mu.RUnlock()

	if !ok {
		return Session{}, errors.New("session not found")
	}

	if time.Now().After(session.ExpiresAt) {
		s.mu.Lock()
		delete(s.sessions, hashedID)
		s.mu.Unlock()
		return Session{}, errors.New("session expired")
	}

	return session, nil
}

func (s *InMemorySessionStore) Delete(ctx context.Context, id string) error {
	hashedID := hashSessionID(id)

	s.mu.Lock()
	delete(s.sessions, hashedID)
	s.mu.Unlock()

	return nil
}

func (s *InMemorySessionStore) Rotate(ctx context.Context, oldID string) (Session, string, error) {
	hashedOldID := hashSessionID(oldID)

	s.mu.Lock()
	oldSession, ok := s.sessions[hashedOldID]
	if !ok {
		s.mu.Unlock()
		return Session{}, "", errors.New("session not found")
	}
	delete(s.sessions, hashedOldID)
	s.mu.Unlock()

	newSession := oldSession
	newSession.ParentID = oldSession.ID
	newSession.IsRotation = true
	newSession.CreatedAt = time.Now()
	newSession.ExpiresAt = time.Now().Add(DefaultSessionDuration)
	if newSession.ExpiresAt.Sub(newSession.CreatedAt) > MaxSessionDuration {
		newSession.ExpiresAt = newSession.CreatedAt.Add(MaxSessionDuration)
	}

	rawID, hashedID, err := generateRawSessionID()
	if err != nil {
		return Session{}, "", err
	}
	newSession.ID = hashedID

	s.mu.Lock()
	s.sessions[hashedID] = newSession
	s.mu.Unlock()

	return newSession, rawID, nil
}

type CookieSessionStore struct {
	cookieName string
	secure     bool
	httpOnly   bool
	sameSite   http.SameSite
	domain     string
	path       string
	store      SessionStore
	maxAge     int
}

func NewCookieSessionStore(cookieName string, secure, httpOnly bool, sameSite http.SameSite) *CookieSessionStore {
	return &CookieSessionStore{
		cookieName: cookieName,
		secure:     secure,
		httpOnly:   httpOnly,
		sameSite:   sameSite,
		path:       "/",
		maxAge:     int(DefaultSessionDuration.Seconds()),
		store:      NewInMemorySessionStore(),
	}
}

func (s *CookieSessionStore) WithStore(store SessionStore) *CookieSessionStore {
	s.store = store
	return s
}

func (s *CookieSessionStore) WithDomain(domain string) *CookieSessionStore {
	s.domain = domain
	return s
}

func (s *CookieSessionStore) WithPath(path string) *CookieSessionStore {
	s.path = path
	return s
}

func (s *CookieSessionStore) WithMaxAge(maxAge int) *CookieSessionStore {
	s.maxAge = maxAge
	return s
}

func (s *CookieSessionStore) buildCookie(value string, expires time.Time) *http.Cookie {
	cookie := &http.Cookie{
		Name:     s.cookieName,
		Value:    value,
		Path:     s.path,
		Domain:   s.domain,
		Expires:  expires,
		MaxAge:   s.maxAge,
		Secure:   s.secure,
		HttpOnly: s.httpOnly,
		SameSite: s.sameSite,
	}
	return cookie
}

func (s *CookieSessionStore) Create(ctx context.Context, session Session) (string, error) {
	return s.store.Create(ctx, session)
}

func (s *CookieSessionStore) Get(ctx context.Context, id string) (Session, error) {
	return s.store.Get(ctx, id)
}

func (s *CookieSessionStore) Delete(ctx context.Context, id string) error {
	return s.store.Delete(ctx, id)
}

func (s *CookieSessionStore) Rotate(ctx context.Context, oldID string) (Session, string, error) {
	return s.store.Rotate(ctx, oldID)
}

func (s *CookieSessionStore) SetSessionCookie(w http.ResponseWriter, rawSessionID string, expires time.Time) {
	http.SetCookie(w, s.buildCookie(rawSessionID, expires))
}

func (s *CookieSessionStore) ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, s.buildCookie("", time.Unix(0, 0)))
}

func (s *CookieSessionStore) GetSecure() bool      { return s.secure }
func (s *CookieSessionStore) GetHTTPOnly() bool     { return s.httpOnly }
func (s *CookieSessionStore) GetSameSite() http.SameSite { return s.sameSite }
