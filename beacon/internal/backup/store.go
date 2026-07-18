package backup

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Store interface {
	Create(ctx context.Context, b Backup) error
	Get(ctx context.Context, id string) (Backup, error)
	List(ctx context.Context, serverID string, limit int) ([]Backup, error)
	UpdateStatus(ctx context.Context, id string, status BackupStatus, errorMsg string) error
	Delete(ctx context.Context, id string) error
}

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(db *sql.DB) *SQLiteStore {
	return &SQLiteStore{db: db}
}

func (s *SQLiteStore) Create(ctx context.Context, b Backup) error {
	const query = `INSERT INTO backups (id, server_id, started_at, completed_at, status, size_bytes, files, duration, adapter, path, error)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := s.db.ExecContext(ctx, query,
		b.ID, b.ServerID, b.StartedAt, b.CompletedAt, b.Status, b.SizeBytes, b.Files,
		b.Duration, b.Adapter, b.Path, b.Error,
	)
	if err != nil {
		return fmt.Errorf("create backup record: %w", err)
	}
	return nil
}

func (s *SQLiteStore) Get(ctx context.Context, id string) (Backup, error) {
	const query = `SELECT id, server_id, started_at, completed_at, status, size_bytes, files, duration, adapter, path, error
		FROM backups WHERE id = ?`
	row := s.db.QueryRowContext(ctx, query, id)
	var b Backup
	var startedAt, completedAt sql.NullTime
	var sizeBytes sql.NullInt64
	var files sql.NullInt64
	var duration sql.NullInt64
	var adapter, path, errMsg sql.NullString
	err := row.Scan(&b.ID, &b.ServerID, &startedAt, &completedAt, &b.Status,
		&sizeBytes, &files, &duration, &adapter, &path, &errMsg,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return Backup{}, fmt.Errorf("backup %s not found", id)
		}
		return Backup{}, fmt.Errorf("get backup %s: %w", id, err)
	}
	if startedAt.Valid {
		b.StartedAt = startedAt.Time
	}
	if completedAt.Valid {
		b.CompletedAt = completedAt.Time
	}
	if sizeBytes.Valid {
		b.SizeBytes = sizeBytes.Int64
	}
	if files.Valid {
		b.Files = int(files.Int64)
	}
	if duration.Valid {
		b.Duration = time.Duration(duration.Int64)
	}
	if adapter.Valid {
		b.Adapter = adapter.String
	}
	if path.Valid {
		b.Path = path.String
	}
	if errMsg.Valid {
		b.Error = errMsg.String
	}
	return b, nil
}

func (s *SQLiteStore) List(ctx context.Context, serverID string, limit int) ([]Backup, error) {
	var query string
	var args []interface{}
	if limit > 0 {
		query = `SELECT id, server_id, started_at, completed_at, status, size_bytes, files, duration, adapter, path, error
			FROM backups WHERE server_id = ? ORDER BY completed_at DESC LIMIT ?`
		args = []interface{}{serverID, limit}
	} else {
		query = `SELECT id, server_id, started_at, completed_at, status, size_bytes, files, duration, adapter, path, error
			FROM backups WHERE server_id = ? ORDER BY completed_at DESC`
		args = []interface{}{serverID}
	}
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list backups for server %s: %w", serverID, err)
	}
	defer rows.Close()

	var backups []Backup
	for rows.Next() {
		var b Backup
		var startedAt, completedAt sql.NullTime
		var sizeBytes sql.NullInt64
		var files sql.NullInt64
		var duration sql.NullInt64
		var adapter, path, errMsg sql.NullString
		err := rows.Scan(&b.ID, &b.ServerID, &startedAt, &completedAt, &b.Status,
			&sizeBytes, &files, &duration, &adapter, &path, &errMsg,
		)
		if err != nil {
			return nil, fmt.Errorf("scan backup row: %w", err)
		}
		if startedAt.Valid {
			b.StartedAt = startedAt.Time
		}
		if completedAt.Valid {
			b.CompletedAt = completedAt.Time
		}
		if sizeBytes.Valid {
			b.SizeBytes = sizeBytes.Int64
		}
		if files.Valid {
			b.Files = int(files.Int64)
		}
		if duration.Valid {
			b.Duration = time.Duration(duration.Int64)
		}
		if adapter.Valid {
			b.Adapter = adapter.String
		}
		if path.Valid {
			b.Path = path.String
		}
		if errMsg.Valid {
			b.Error = errMsg.String
		}
		backups = append(backups, b)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate backup rows: %w", err)
	}
	return backups, nil
}

func (s *SQLiteStore) UpdateStatus(ctx context.Context, id string, status BackupStatus, errorMsg string) error {
	const query = `UPDATE backups SET status = ?, error = ? WHERE id = ?`
	result, err := s.db.ExecContext(ctx, query, status, errorMsg, id)
	if err != nil {
		return fmt.Errorf("update backup %s status: %w", id, err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check rows affected for backup %s: %w", id, err)
	}
	if affected == 0 {
		return fmt.Errorf("backup %s not found", id)
	}
	return nil
}

func (s *SQLiteStore) Delete(ctx context.Context, id string) error {
	const query = `DELETE FROM backups WHERE id = ?`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete backup %s: %w", id, err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check rows affected for backup %s: %w", id, err)
	}
	if affected == 0 {
		return fmt.Errorf("backup %s not found", id)
	}
	return nil
}
