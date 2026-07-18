//go:build webauthn

package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

const (
	SessionDataLength = 64
)

type WebAuthnUser struct {
	ID          []byte
	Name        string
	DisplayName string
	Credentials []webauthn.Credential
}

func (u *WebAuthnUser) WebAuthnID() []byte {
	return u.ID
}

func (u *WebAuthnUser) WebAuthnName() string {
	return u.Name
}

func (u *WebAuthnUser) WebAuthnDisplayName() string {
	return u.DisplayName
}

func (u *WebAuthnUser) WebAuthnCredentials() []webauthn.Credential {
	return u.Credentials
}

type WebAuthn struct {
	rpID             string
	rpOrigins        []string
	rpDisplayName    string
	webauthn         *webauthn.WebAuthn
	sessionStorage   WebAuthnSessionStorage
}

type WebAuthnSessionStorage interface {
	StoreRegistrationSession(userID string, sessionData []byte) error
	GetRegistrationSession(userID string) ([]byte, error)
	DeleteRegistrationSession(userID string) error
	StoreLoginSession(userID string, sessionData []byte) error
	GetLoginSession(userID string) ([]byte, error)
	DeleteLoginSession(userID string) error
}

type InMemoryWebAuthnSessionStore struct {
	mu            sync.RWMutex
	regSessions   map[string][]byte
	loginSessions map[string][]byte
}

func NewInMemoryWebAuthnSessionStore() *InMemoryWebAuthnSessionStore {
	return &InMemoryWebAuthnSessionStore{
		regSessions:   make(map[string][]byte),
		loginSessions: make(map[string][]byte),
	}
}

func (s *InMemoryWebAuthnSessionStore) StoreRegistrationSession(userID string, sessionData []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.regSessions[userID] = sessionData
	return nil
}

func (s *InMemoryWebAuthnSessionStore) GetRegistrationSession(userID string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	data, ok := s.regSessions[userID]
	if !ok {
		return nil, fmt.Errorf("registration session not found for user: %s", userID)
	}
	return data, nil
}

func (s *InMemoryWebAuthnSessionStore) DeleteRegistrationSession(userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.regSessions, userID)
	return nil
}

func (s *InMemoryWebAuthnSessionStore) StoreLoginSession(userID string, sessionData []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.loginSessions[userID] = sessionData
	return nil
}

func (s *InMemoryWebAuthnSessionStore) GetLoginSession(userID string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	data, ok := s.loginSessions[userID]
	if !ok {
		return nil, fmt.Errorf("login session not found for user: %s", userID)
	}
	return data, nil
}

func (s *InMemoryWebAuthnSessionStore) DeleteLoginSession(userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.loginSessions, userID)
	return nil
}

func NewWebAuthn(rpID string, rpOrigins []string, rpDisplayName string) (*WebAuthn, error) {
	if rpID == "" {
		return nil, fmt.Errorf("rpID is required")
	}
	if len(rpOrigins) == 0 {
		return nil, fmt.Errorf("at least one rpOrigin is required")
	}
	if rpDisplayName == "" {
		return nil, fmt.Errorf("rpDisplayName is required")
	}

	wconfig := &webauthn.Config{
		RPDisplayName:        rpDisplayName,
		RPID:                 rpID,
		RPOrigins:            rpOrigins,
		AttestationPreference: protocol.PreferenceDirect,
		Timeouts: webauthn.TimeoutsConfig{
			RegistrationTimeout: 5 * time.Minute,
			LoginTimeout:        5 * time.Minute,
		},
	}

	wa, err := webauthn.New(wconfig)
	if err != nil {
		return nil, err
	}

	return &WebAuthn{
		rpID:           rpID,
		rpOrigins:      rpOrigins,
		rpDisplayName:  rpDisplayName,
		webauthn:       wa,
		sessionStorage: NewInMemoryWebAuthnSessionStore(),
	}, nil
}

func (w *WebAuthn) WithSessionStorage(storage WebAuthnSessionStorage) *WebAuthn {
	w.sessionStorage = storage
	return w
}

func (w *WebAuthn) GenerateUserID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate user ID: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func NewWebAuthnUser(id, name, displayName string) *WebAuthnUser {
	userID := []byte(id)
	return &WebAuthnUser{
		ID:          userID,
		Name:        name,
		DisplayName: displayName,
	}
}

func (w *WebAuthn) BeginRegistration(user *WebAuthnUser, opts ...webauthn.RegistrationOption) (*protocol.CredentialCreation, *webauthn.SessionData, error) {
	cc, sessionData, err := w.webauthn.BeginRegistration(user, opts...)
	if err != nil {
		return nil, nil, err
	}

	if w.sessionStorage != nil {
		sessionBytes, err := json.Marshal(sessionData)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to marshal registration session data: %w", err)
		}
		if err := w.sessionStorage.StoreRegistrationSession(string(user.WebAuthnID()), sessionBytes); err != nil {
			return nil, nil, fmt.Errorf("failed to store registration session: %w", err)
		}
	}

	return cc, sessionData, nil
}

func (w *WebAuthn) FinishRegistration(user *WebAuthnUser, r *http.Request) (*webauthn.Credential, error) {
	var sessionData webauthn.SessionData

	if w.sessionStorage != nil {
		sessionBytes, err := w.sessionStorage.GetRegistrationSession(string(user.WebAuthnID()))
		if err != nil {
			return nil, fmt.Errorf("failed to get registration session: %w", err)
		}
		if err := json.Unmarshal(sessionBytes, &sessionData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal registration session data: %w", err)
		}
	}

	credential, err := w.webauthn.FinishRegistration(user, sessionData, r)
	if err != nil {
		return nil, fmt.Errorf("failed to finish registration: %w", err)
	}

	if w.sessionStorage != nil {
		if err := w.sessionStorage.DeleteRegistrationSession(string(user.WebAuthnID())); err != nil {
			return nil, err
		}
	}

	user.Credentials = append(user.Credentials, *credential)

	return credential, nil
}

func (w *WebAuthn) BeginLogin(user *WebAuthnUser, opts ...webauthn.LoginOption) (*protocol.CredentialAssertion, *webauthn.SessionData, error) {
	assertion, sessionData, err := w.webauthn.BeginLogin(user, opts...)
	if err != nil {
		return nil, nil, err
	}

	if w.sessionStorage != nil {
		sessionBytes, err := json.Marshal(sessionData)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to marshal login session data: %w", err)
		}
		if err := w.sessionStorage.StoreLoginSession(string(user.WebAuthnID()), sessionBytes); err != nil {
			return nil, nil, fmt.Errorf("failed to store login session: %w", err)
		}
	}

	return assertion, sessionData, nil
}

func (w *WebAuthn) FinishLogin(user *WebAuthnUser, r *http.Request) (*webauthn.Credential, error) {
	var sessionData webauthn.SessionData

	if w.sessionStorage != nil {
		sessionBytes, err := w.sessionStorage.GetLoginSession(string(user.WebAuthnID()))
		if err != nil {
			return nil, fmt.Errorf("failed to get login session: %w", err)
		}
		if err := json.Unmarshal(sessionBytes, &sessionData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal login session data: %w", err)
		}
	}

	credential, err := w.webauthn.FinishLogin(user, sessionData, r)
	if err != nil {
		return nil, fmt.Errorf("failed to finish login: %w", err)
	}

	if w.sessionStorage != nil {
		if err := w.sessionStorage.DeleteLoginSession(string(user.WebAuthnID())); err != nil {
			return nil, err
		}
	}

	return credential, nil
}
