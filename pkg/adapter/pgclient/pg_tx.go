package pgclient

import (
	"context"
	"database/sql"
	"errors"

	"go.uber.org/zap"
)

type pgTx struct {
	ctx           context.Context
	logger        Logger
	cfg           *Config
	tx            *sql.Tx
	hookFactories hookFactories
}

func (t *pgTx) Begin(context.Context) (DB, error) {
	return t, nil
}

func (t *pgTx) Commit() error {
	if err := t.tx.Commit(); err != nil && !errors.Is(err, sql.ErrTxDone) {
		t.logger.Error("cannot commit transaction", zap.Error(err))
		return err
	}

	return nil
}

func (t *pgTx) Rollback(context.Context) error {
	if err := t.tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
		t.logger.Error("cannot rollback transaction", zap.Error(err))
		return err
	}

	return nil
}

func (t *pgTx) DB() *sql.DB {
	return nil
}

func (t *pgTx) Close() error {
	return errors.New("unsupported operation for transaction")
}

func (t *pgTx) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	hooks := t.hookFactories.createHooks()

	for _, hook := range hooks {
		hook.Before(ctx)
	}

	result, err := t.tx.ExecContext(ctx, query, args...)

	for _, hook := range hooks {
		hook.After(ctx, query, args...)
	}

	return result, err
}

func (t *pgTx) Exec(query string, args ...any) (sql.Result, error) {
	return t.ExecContext(t.ctx, query, args...)
}

func (t *pgTx) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	hooks := t.hookFactories.createHooks()

	for _, hook := range hooks {
		hook.Before(ctx)
	}

	result, err := t.tx.QueryContext(ctx, query, args...)

	for _, hook := range hooks {
		hook.After(ctx, query, args...)
	}

	return result, err
}

func (t *pgTx) Query(query string, args ...any) (*sql.Rows, error) {
	return t.QueryContext(t.ctx, query, args...)
}

func (t *pgTx) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	hooks := t.hookFactories.createHooks()

	for _, hook := range hooks {
		hook.Before(ctx)
	}

	row := t.tx.QueryRowContext(ctx, query, args...)

	for _, hook := range hooks {
		hook.After(ctx, query, args...)
	}

	return row
}

func (t *pgTx) QueryRow(query string, args ...any) *sql.Row {
	return t.QueryRowContext(t.ctx, query, args...)
}
