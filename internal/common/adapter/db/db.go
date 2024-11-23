package db

import (
	"context"
	"database/sql"
)

type DB interface {
	Begin(context.Context) (DB, error)
	Commit() error
	Rollback() error
	Close() error

	ExecContext(context.Context, string, ...any) (sql.Result, error)
	Exec(query string, args ...interface{}) (sql.Result, error)

	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)

	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
	QueryRow(string, ...interface{}) *sql.Row
}
