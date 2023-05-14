package db

import (
	"context"
	"database/sql"

	"github.com/vizitiuRoman/go-grpc-boilerplate/internal/common/adapter/log"
)

type pgTx struct {
	ctx    context.Context
	tx     *sql.Tx
	logger log.Logger
}

func (t *pgTx) Begin(context.Context) (DB, error) {
	return t, nil
}

func (t *pgTx) Commit() error {
	t.logger.Debug("commit transaction")
	return t.tx.Commit()
}

func (t *pgTx) Rollback() error {
	return t.tx.Rollback()
}

func (t *pgTx) Close() error {
	return nil
}
func (t *pgTx) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return t.tx.ExecContext(ctx, query, args...)
}

func (t *pgTx) Exec(query string, args ...any) (sql.Result, error) {
	return t.ExecContext(t.ctx, query, args...)
}

func (t *pgTx) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return t.tx.QueryContext(ctx, query, args...)
}

func (t *pgTx) Query(query string, args ...any) (*sql.Rows, error) {
	return t.QueryContext(t.ctx, query, args...)
}

func (t *pgTx) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return t.tx.QueryRowContext(ctx, query, args...)
}

func (t *pgTx) QueryRow(query string, args ...any) *sql.Row {
	return t.tx.QueryRowContext(t.ctx, query, args...)
}
