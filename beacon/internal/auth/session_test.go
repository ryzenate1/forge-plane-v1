package auth

import (
	"context"
	"net/http"
	"testing"
	"time"
)

func TestInMemorySessionStore_Create(t *testing.T) {
	store := NewInMemorySessionStore()
	session := Session{
		UserID:    "test-user",
		Scopes:    Scopes{ScopeServerRead, ScopeServerWrite},
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	rawID, err := store.Create(context.Background(), session)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if rawID == "" {
		t.Error("Create() returned empty rawID")
	}
}

func TestInMemorySessionStore_Get(t *testing.T) {
	store := NewInMemorySessionStore()
	session := Session{
		UserID:    "test-user",
		Scopes:    Scopes{ScopeServerRead, ScopeServerWrite},
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	rawID, err := store.Create(context.Background(), session)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	retrieved, err := store.Get(context.Background(), rawID)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if retrieved.UserID != session.UserID {
		t.Errorf("Get() UserID = %v, want %v", retrieved.UserID, session.UserID)
	}
}

func TestInMemorySessionStore_Delete(t *testing.T) {
	store := NewInMemorySessionStore()
	session := Session{
		UserID:    "test-user",
		Scopes:    Scopes{ScopeServerRead, ScopeServerWrite},
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	rawID, err := store.Create(context.Background(), session)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	err = store.Delete(context.Background(), rawID)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	_, err = store.Get(context.Background(), rawID)
	if err == nil {
		t.Error("Get() did not return error for deleted session")
	}
}

func TestCookieSessionStore_CreateGetDelete(t *testing.T) {
	store := NewCookieSessionStore("session", true, true, http.SameSiteLaxMode)
	session := Session{
		UserID:    "test-user",
		Scopes:    Scopes{ScopeServerRead, ScopeServerWrite},
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	rawID, err := store.Create(context.Background(), session)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if rawID == "" {
		t.Error("Create() returned empty rawID")
	}

	retrieved, err := store.Get(context.Background(), rawID)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if retrieved.UserID != session.UserID {
		t.Errorf("Get() UserID = %v, want %v", retrieved.UserID, session.UserID)
	}

	err = store.Delete(context.Background(), rawID)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	_, err = store.Get(context.Background(), rawID)
	if err == nil {
		t.Error("Get() did not return error for deleted session")
	}
}
