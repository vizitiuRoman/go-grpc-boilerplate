package pgclient

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

func NewPool(ctx context.Context, cfg *Config, logger Logger, opts ...Option) (DB, error) {
	db, err := sql.Open("pgx", cfg.DSN)
	if err != nil {
		logger.Error("cannot to open pgclient", zap.Error(err))
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxIdleTime(time.Second * time.Duration(cfg.ConnMaxIdleTime))

	if err := db.Ping(); err != nil {
		return nil, err
	}

	res := &pgDB{
		ctx:    ctx,
		logger: logger,
		cfg:    cfg,
		db:     db,
	}

	for _, opt := range opts {
		opt(res)
	}

	return res, nil
}

type pgDB struct {
	ctx           context.Context
	logger        Logger
	cfg           *Config
	db            *sql.DB
	hookFactories hookFactories
}

func (d *pgDB) Begin(ctx context.Context) (DB, error) {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &pgTx{
		ctx:           ctx,
		logger:        d.logger,
		cfg:           d.cfg,
		tx:            tx,
		hookFactories: d.hookFactories,
	}, nil
}

func (d *pgDB) Commit() error {
	return nil
}

func (d *pgDB) Rollback(context.Context) error { return nil }

func (d *pgDB) Close() error {
	if err := d.db.Close(); err != nil {
		d.logger.Error("cannot close pgclient", zap.Error(err))
		return err
	}

	d.logger.Debug("pgclient pool closed")
	return nil
}

func (d *pgDB) DB() *sql.DB {
	return d.db
}

func (d *pgDB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	hooks := d.hookFactories.createHooks()

	for _, hook := range hooks {
		hook.Before(ctx)
	}

	result, err := d.db.ExecContext(ctx, query, args...)

	for _, hook := range hooks {
		hook.After(ctx, query, args...)
	}

	return result, err
}

func (d *pgDB) Exec(query string, args ...any) (sql.Result, error) {
	return d.ExecContext(d.ctx, query, args...)
}

func (d *pgDB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	hooks := d.hookFactories.createHooks()

	for _, hook := range hooks {
		hook.Before(ctx)
	}

	result, err := d.db.QueryContext(ctx, query, args...)

	for _, hook := range hooks {
		hook.After(ctx, query, args...)
	}

	return result, err
}

func (d *pgDB) Query(query string, args ...any) (*sql.Rows, error) {
	return d.QueryContext(d.ctx, query, args...)
}

func (d *pgDB) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	hooks := d.hookFactories.createHooks()

	for _, hook := range hooks {
		hook.Before(ctx)
	}

	row := d.db.QueryRowContext(ctx, query, args...)

	for _, hook := range hooks {
		hook.After(ctx, query, args...)
	}

	return row
}

func (d *pgDB) QueryRow(query string, args ...any) *sql.Row {
	return d.QueryRowContext(d.ctx, query, args...)
}
