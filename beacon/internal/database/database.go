package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/glebarez/go-sqlite"
)

type Database interface {
	Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row
	Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	BeginTx(ctx context.Context) (*sql.Tx, error)
	Ping(ctx context.Context) error
	Close() error
}

type SQLiteDatabase struct {
	db *sql.DB
}

func NewSQLiteDatabase(dsn string) (*SQLiteDatabase, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(3)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	d := &SQLiteDatabase{db: db}
	if err := d.Ping(context.Background()); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return d, nil
}

func (s *SQLiteDatabase) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return s.db.QueryContext(ctx, query, args...)
}

func (s *SQLiteDatabase) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return s.db.QueryRowContext(ctx, query, args...)
}

func (s *SQLiteDatabase) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return s.db.ExecContext(ctx, query, args...)
}

func (s *SQLiteDatabase) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return s.db.BeginTx(ctx, nil)
}

func (s *SQLiteDatabase) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

func (s *SQLiteDatabase) Close() error {
	return s.db.Close()
}
