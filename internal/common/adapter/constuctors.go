package adapter

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/vizitiuRoman/go-grpc-boilerplate/internal/common/adapter/config"
	"github.com/vizitiuRoman/go-grpc-boilerplate/internal/common/adapter/db"
	"github.com/vizitiuRoman/go-grpc-boilerplate/internal/common/adapter/log"
	"github.com/vizitiuRoman/go-grpc-boilerplate/internal/common/adapter/meta"
)

var Constructors = fx.Provide(
	config.NewConfig,

	NewFxLogger,

	fx.Annotate(NewFxPgPool, fx.OnStop(func(db db.DB) error { return db.Close() })),
)

func NewFxLogger(cfg *config.Config) (log.Logger, error) {
	return log.NewLogger(cfg.Logger, zap.String(meta.AppVersionKey, cfg.Version+" "+cfg.BuildDate))
}

func NewFxPgPool(ctx context.Context, logger log.Logger, cfg *config.Config) (db.DB, error) {
	return db.NewPool(ctx, logger, cfg.DB)
}
