package backup

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })
	_, err = db.Exec(`
        CREATE TABLE backups (
            id TEXT PRIMARY KEY,
            server_id TEXT,
            started_at DATETIME,
            completed_at DATETIME,
            status TEXT,
            size_bytes INTEGER,
            files INTEGER,
            duration INTEGER,
            adapter TEXT,
            path TEXT,
            error TEXT
        );
    `)
	require.NoError(t, err)
	return db
}

func TestVerifyBackup(t *testing.T) {
	store := NewSQLiteStore(setupTestDB(t))
	mock := NewMockBackup()

	backupID := "test-backup.zip"
	mock.AddBackup("server-1", backupID, []byte("test-data"), time.Now())

	testBackup := Backup{
		ID:          backupID,
		ServerID:    "server-1",
		StartedAt:   time.Now(),
		CompletedAt: time.Now().Add(time.Minute),
		Status:      BackupStatusCompleted,
		SizeBytes:   1024,
		Files:       5,
		Duration:    time.Minute,
		Adapter:     "local",
		Path:        "/backups/test",
	}

	err := store.Create(context.Background(), testBackup)
	require.NoError(t, err)

	result, err := VerifyBackup(context.Background(), mock, "server-1", backupID)
	assert.NoError(t, err)
	assert.Equal(t, BackupStatusCompleted, result.Status)
	assert.Equal(t, backupID, result.BackupID)
	assert.NotEmpty(t, result.Checksum)

	updatedBackup, err := store.Get(context.Background(), backupID)
	assert.NoError(t, err)
	assert.Equal(t, BackupStatusCompleted, updatedBackup.Status)
}
