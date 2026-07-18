package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAuthMiddleware(t *testing.T) {
	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Wrap the test handler with the middleware
	wrappedHandler := AuthMiddleware(testHandler)

	// Create a test request
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	// Serve the request
	wrappedHandler.ServeHTTP(w, req)

	// Check the response
	if w.Code != http.StatusOK {
		t.Errorf("AuthMiddleware() status code = %v, want %v", w.Code, http.StatusOK)
	}
}

func TestRequireScopes(t *testing.T) {
	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Wrap the test handler with the middleware
	wrappedHandler := RequireScopes(ScopeServerRead, ScopeServerWrite)(testHandler)

	// Create a test request with a valid session in context
	req := httptest.NewRequest("GET", "/", nil)
	session := Session{
		Scopes:    Scopes{ScopeServerRead, ScopeServerWrite},
		ExpiresAt: time.Now().Add(time.Hour),
	}
	ctx := context.WithValue(req.Context(), ContextKeySession, session)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// Serve the request
	wrappedHandler.ServeHTTP(w, req)

	// Check the response
	if w.Code != http.StatusOK {
		t.Errorf("RequireScopes() status code = %v, want %v", w.Code, http.StatusOK)
	}
}
