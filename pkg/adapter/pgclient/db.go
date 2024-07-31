package pgclient

import (
	"context"
	"database/sql"
)

type DB interface {
	Begin(context.Context) (DB, error)
	Commit() error
	Rollback(context.Context) error

	DB() *sql.DB
	Close() error

	ExecContext(context.Context, string, ...any) (sql.Result, error)
	Exec(query string, args ...any) (sql.Result, error)

	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	Query(query string, args ...any) (*sql.Rows, error)

	QueryRowContext(context.Context, string, ...any) *sql.Row
	QueryRow(string, ...any) *sql.Row
}
