package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/vizitiuRoman/go-grpc-boilerplate/internal/common/adapter/log"
	"go.uber.org/zap"

	// init sql driver
	_ "github.com/lib/pq"
)

func NewPool(ctx context.Context, logger log.Logger, cnf *Config) (DB, error) {
	logger = logger.WithComponent(ctx, "db")

	db, err := sql.Open("postgres", cnf.DSN)
	if err != nil {
		logger.Error("failed to open db", zap.Error(err))
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(30)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(time.Minute * 4)

	return &pgDB{
		ctx:    ctx,
		db:     db,
		logger: logger,
	}, nil
}

type pgDB struct {
	ctx    context.Context
	db     *sql.DB
	logger log.Logger
}

func (d *pgDB) Begin(ctx context.Context) (DB, error) {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	d.logger.Debug("beginning transaction")

	return &pgTx{
		logger: d.logger,
		ctx:    ctx,
		tx:     tx,
	}, nil
}

func (d *pgDB) Commit() error {
	return nil
}

func (d *pgDB) Rollback() error {
	return nil
}

func (d *pgDB) Close() error {
	return d.db.Close()
}

func (d *pgDB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return d.db.ExecContext(ctx, query, args...)
}

func (d *pgDB) Exec(query string, args ...any) (sql.Result, error) {
	return d.ExecContext(d.ctx, query, args...)
}

func (d *pgDB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return d.db.QueryContext(ctx, query, args...)
}

func (d *pgDB) Query(query string, args ...any) (*sql.Rows, error) {
	return d.QueryContext(d.ctx, query, args...)
}

func (d *pgDB) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return d.db.QueryRowContext(ctx, query, args...)
}

func (d *pgDB) QueryRow(query string, args ...any) *sql.Row {
	return d.db.QueryRowContext(d.ctx, query, args...)
}
