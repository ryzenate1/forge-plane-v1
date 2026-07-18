package auth

import (
	"context"
	"crypto/subtle"
	"errors"
	"net/http"
	"strings"
	"time"
)

type contextKey string

const (
	ContextKeySession contextKey = "auth:session"
)

var (
	ErrNoAuthHeader     = errors.New("missing authorization header")
	ErrInvalidAuthHeader = errors.New("invalid authorization header format")
	ErrSessionExpired   = errors.New("session expired")
	ErrInsufficientScope = errors.New("insufficient permissions")
	ErrSessionNotFound  = errors.New("session not found")
)

var sessionStore SessionStore

func SetSessionStore(store SessionStore) {
	sessionStore = store
}

func GetSessionFromContext(ctx context.Context) (Session, bool) {
	s, ok := ctx.Value(ContextKeySession).(Session)
	return s, ok
}

func extractBearerToken(r *http.Request) (string, error) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		cookie, err := r.Cookie("session_token")
		if err == nil && cookie.Value != "" {
			return cookie.Value, nil
		}
		return "", ErrNoAuthHeader
	}

	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		parts = strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return "", ErrInvalidAuthHeader
		}
	}

	return parts[1], nil
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if sessionStore == nil {
			next.ServeHTTP(w, r)
			return
		}

		token, err := extractBearerToken(r)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		session, err := sessionStore.Get(r.Context(), token)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		if time.Now().After(session.ExpiresAt) {
			sessionStore.Delete(r.Context(), token)
			http.Error(w, "session expired", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ContextKeySession, session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequireScopes(required ...Scope) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, ok := GetSessionFromContext(r.Context())
			if !ok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			if !session.Scopes.HasAll(required...) {
				http.Error(w, ErrInsufficientScope.Error(), http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func ValidateSession(r *http.Request) (Session, error) {
	if sessionStore == nil {
		return Session{}, errors.New("session store not configured")
	}

	token, err := extractBearerToken(r)
	if err != nil {
		return Session{}, err
	}

	session, err := sessionStore.Get(r.Context(), token)
	if err != nil {
		return Session{}, ErrSessionNotFound
	}

	if time.Now().After(session.ExpiresAt) {
		sessionStore.Delete(r.Context(), token)
		return Session{}, ErrSessionExpired
	}

	return session, nil
}

func CheckScopes(r *http.Request, required ...Scope) error {
	session, err := ValidateSession(r)
	if err != nil {
		return err
	}

	if !session.Scopes.HasAll(required...) {
		return ErrInsufficientScope
	}

	return nil
}

func constantTimeCompare(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
