package migration

import (
	"context"
	"errors"
	"testing"

	gpruntime "gamepanel/forge/internal/runtime"
)

func TestExecuteMigrationReturnsTypedNotImplementedWithoutStore(t *testing.T) {
	service := New(nil, nil, nil, nil, nil)

	_, err := service.ExecuteMigration(context.Background(), "migration-1")
	if err == nil {
		t.Fatal("expected migration execution to be unavailable")
	}
	var notImplemented *NotImplementedError
	if !errors.As(err, &notImplemented) {
		t.Fatalf("expected *NotImplementedError, got %T", err)
	}
	if notImplemented.MigrationID != "migration-1" {
		t.Fatalf("unexpected migration id %q", notImplemented.MigrationID)
	}
	if !errors.Is(err, gpruntime.ErrNotImplemented) {
		t.Fatalf("expected error to wrap runtime.ErrNotImplemented: %v", err)
	}
}

func TestPrepareMigrationReturnsTypedNotImplementedWithoutStore(t *testing.T) {
	service := New(nil, nil, nil, nil, nil)

	_, err := service.PrepareMigration(context.Background(), "migration-1")
	var notImplemented *NotImplementedError
	if !errors.As(err, &notImplemented) {
		t.Fatalf("expected *NotImplementedError, got %T (%v)", err, err)
	}
}
