package recovery

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"gamepanel/forge/internal/daemon"
	"gamepanel/forge/internal/store"
)

type recordingMigrationExecutor struct {
	calls      []string
	prepareErr error
	executeErr error
}

func (e *recordingMigrationExecutor) PrepareMigration(_ context.Context, id string) (store.Migration, error) {
	e.calls = append(e.calls, "prepare:"+id)
	return store.Migration{ID: id}, e.prepareErr
}

func (e *recordingMigrationExecutor) ExecuteMigration(_ context.Context, id string) (store.Migration, error) {
	e.calls = append(e.calls, "execute:"+id)
	return store.Migration{ID: id}, e.executeErr
}

func (e *recordingMigrationExecutor) GetMigration(_ context.Context, id string) (store.Migration, error) {
	return store.Migration{ID: id}, nil
}

func (e *recordingMigrationExecutor) CancelMigration(_ context.Context, id string) (store.Migration, error) {
	return store.Migration{ID: id}, nil
}

func (e *recordingMigrationExecutor) MarkFailed(_ context.Context, id, _ string) (store.Migration, error) {
	return store.Migration{ID: id}, nil
}

var _ MigrationExecutor = (*recordingMigrationExecutor)(nil)

type recordingBackupRestoreExecutor struct {
	calls []store.RecoveryItem
	err   error
}

func (e *recordingBackupRestoreExecutor) VerifyAndRestore(_ context.Context, item store.RecoveryItem) error {
	e.calls = append(e.calls, item)
	return e.err
}

var _ BackupRestoreExecutor = (*recordingBackupRestoreExecutor)(nil)

func TestNewWithMigrationExecutorRegistersExecutor(t *testing.T) {
	executor := &recordingMigrationExecutor{}
	coordinator := NewWithMigrationExecutor(nil, nil, nil, executor)

	if coordinator.migrationExecutor() != executor {
		t.Fatal("migration executor was not registered")
	}
}

func TestBackupRestoreExecutorRegistersWithoutMigrationExecutor(t *testing.T) {
	restore := &recordingBackupRestoreExecutor{}
	coordinator := New(nil, nil, nil)
	coordinator.SetBackupRestoreExecutor(restore)

	if coordinator.backupRestoreExecutor() != restore {
		t.Fatal("backup restore executor was not registered")
	}
	if coordinator.migrationExecutor() != nil {
		t.Fatal("backup recovery must not require a live migration executor")
	}
}

func TestDaemonBackupRestoreRejectsUnverifiedSource(t *testing.T) {
	executor := &DaemonBackupRestoreExecutor{store: &store.Store{}, daemon: daemon.NewClient()}
	err := executor.VerifyAndRestore(context.Background(), store.RecoveryItem{
		ID: "item-1", ServerID: "server-1", TargetNodeID: "target-1",
	})
	if err == nil || err.Error() != "recovery item has no verified backup restore source" {
		t.Fatalf("VerifyAndRestore error = %v, want missing-source error", err)
	}
}

func TestExecuteItemPreparesThenExecutesMigration(t *testing.T) {
	executor := &recordingMigrationExecutor{}
	coordinator := &Coordinator{}
	item := store.RecoveryItem{ID: "item-1", MigrationID: "migration-1"}

	if err := coordinator.executeItem(context.Background(), executor, item); err != nil {
		t.Fatalf("executeItem returned an error: %v", err)
	}
	if want := []string{"prepare:migration-1", "execute:migration-1"}; !reflect.DeepEqual(executor.calls, want) {
		t.Fatalf("executor calls = %v, want %v", executor.calls, want)
	}
}

func TestExecuteItemDoesNotExecuteWhenPreparationFails(t *testing.T) {
	prepareErr := errors.New("target allocation unavailable")
	executor := &recordingMigrationExecutor{prepareErr: prepareErr}
	coordinator := &Coordinator{}
	item := store.RecoveryItem{ID: "item-1", MigrationID: "migration-1"}

	err := coordinator.executeItem(context.Background(), executor, item)
	if !errors.Is(err, prepareErr) {
		t.Fatalf("executeItem error = %v, want wrapped %v", err, prepareErr)
	}
	if want := []string{"prepare:migration-1"}; !reflect.DeepEqual(executor.calls, want) {
		t.Fatalf("executor calls = %v, want %v", executor.calls, want)
	}
}

func TestExecuteItemRequiresMigrationRecord(t *testing.T) {
	executor := &recordingMigrationExecutor{}
	coordinator := &Coordinator{}

	err := coordinator.executeItem(context.Background(), executor, store.RecoveryItem{ID: "item-1"})
	if err == nil || err.Error() != "recovery item has no migration" {
		t.Fatalf("executeItem error = %v, want missing-migration error", err)
	}
	if len(executor.calls) != 0 {
		t.Fatalf("executor was called for an item without a migration: %v", executor.calls)
	}
}
