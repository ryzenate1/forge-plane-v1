package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Migrator struct {
	db     Database
	dir    string
	logger *log.Logger
}

func NewMigrator(db Database, dir string, logger *log.Logger) *Migrator {
	return &Migrator{
		db:     db,
		dir:    dir,
		logger: logger,
	}
}

func (m *Migrator) Migrate(ctx context.Context) error {
	if err := m.db.Ping(ctx); err != nil {
		return fmt.Errorf("database unreachable: %w", err)
	}

	_, err := m.db.Exec(ctx,
		`CREATE TABLE IF NOT EXISTS migrations (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL UNIQUE,
            applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
        );`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	entries, err := os.ReadDir(m.dir)
	if err != nil {
		return fmt.Errorf("failed to read migration directory: %w", err)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	rows, err := m.db.Query(ctx, "SELECT name FROM migrations ORDER BY name")
	if err != nil {
		return fmt.Errorf("failed to query applied migrations: %w", err)
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return fmt.Errorf("failed to scan applied migration: %w", err)
		}
		applied[name] = true
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating applied migrations: %w", err)
	}

	for _, entry := range entries {
		select {
		case <-ctx.Done():
			return fmt.Errorf("migration aborted: %w", ctx.Err())
		default:
		}

		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		if applied[entry.Name()] {
			m.logger.Printf("Migration %s already applied, skipping", entry.Name())
			continue
		}

		m.logger.Printf("Applying migration %s", entry.Name())

		content, err := os.ReadFile(filepath.Join(m.dir, entry.Name()))
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", entry.Name(), err)
		}

		tx, err := m.db.BeginTx(ctx)
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}

		if _, err := tx.ExecContext(ctx, string(content)); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %s: %w", entry.Name(), err)
		}

		if _, err := tx.ExecContext(ctx, "INSERT INTO migrations (name) VALUES (?)", entry.Name()); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %s: %w", entry.Name(), err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", entry.Name(), err)
		}
	}

	m.logger.Println("Migrations applied successfully")
	return nil
}
