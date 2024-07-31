package adapter

import (
	"context"

	log "github.com/vizitiuRoman/go-grpc-boilerplate/pkg/adapter/logger"
	"github.com/vizitiuRoman/go-grpc-boilerplate/pkg/adapter/pgclient"
	"go.uber.org/fx"

	"github.com/vizitiuRoman/go-grpc-boilerplate/internal/common/adapter/config"
)

var Constructors = fx.Provide(
	config.NewConfig,

	NewFxLogger,

	fx.Annotate(NewFxPgPool, fx.OnStop(func(db pgclient.DB) error { return db.Close() })),
)

func NewFxLogger(cfg *config.Config) (log.Logger, error) {
	return log.NewLogger(cfg.Logger)
}

func NewFxPgPool(lf fx.Lifecycle, ctx context.Context, logger log.Logger, cfg *config.Config) (pgclient.DB, error) {
	pool, err := pgclient.NewPool(ctx, cfg.DB, logger)
	lf.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return pool.Close()
		},
	})

	return pool, err
}
