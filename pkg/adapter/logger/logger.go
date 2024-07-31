package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)

	WithComponent(context.Context, string) Logger
	WithMethod(context.Context, string) Logger

	AddContext(ctx context.Context) Logger
	With(fields ...zap.Field) Logger
	UnZap() *zap.Logger
	Zap(*zap.Logger)
}

type logger struct {
	*zap.Logger
	cfg            *Config
	contextEncoder ContextEncoder
}

func NewLogger(cfg *Config, opts ...Option) (Logger, error) {
	zapConfig := defaultJSONConfig(cfg.Level)

	zapConfig.Encoding = cfg.Encoding

	var le zapcore.LevelEncoder
	_ = le.UnmarshalText([]byte(cfg.LevelEncoder))
	zapConfig.EncoderConfig.EncodeLevel = le

	baseLogger, err := zapConfig.Build()
	if err != nil {
		return nil, err
	}
	defer func() { _ = baseLogger.Sync() }()

	if cfg.CallerSkip != nil {
		baseLogger.WithOptions(zap.AddCallerSkip(*cfg.CallerSkip))
	}

	l := &logger{cfg: cfg, Logger: baseLogger}

	for _, opt := range opts {
		opt(l)
	}

	return l, nil
}

func defaultJSONConfig(lvl string) *zap.Config {
	return &zap.Config{
		Level:            unmarshalLevel(lvl),
		DisableCaller:    false,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
		Encoding:         "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "msg",
			LevelKey:       "level",
			TimeKey:        "ts_date",
			NameKey:        "logger",
			CallerKey:      "caller",
			StacktraceKey:  "stacktrace",
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,
		},
	}
}

func unmarshalLevel(lvl string) zap.AtomicLevel {
	level := zap.NewAtomicLevel()
	err := level.UnmarshalText([]byte(lvl))
	if err != nil || lvl == "" {
		level.SetLevel(zap.DebugLevel)
	}

	return level
}

func (l *logger) AddContext(ctx context.Context) Logger {
	if l.contextEncoder != nil {
		return l.contextEncoder(ctx, l)
	}

	return l
}

func (l *logger) WithComponent(ctx context.Context, componentName string) Logger {
	const key = "go.component"

	return l.AddContext(ctx).With(zap.String(key, componentName))
}

func (l *logger) WithMethod(ctx context.Context, methodName string) Logger {
	const key = "go.method"

	return l.AddContext(ctx).With(zap.String(key, methodName))
}

func (l *logger) With(fields ...zap.Field) Logger {
	return &logger{
		cfg:    l.cfg,
		Logger: l.Logger.With(fields...),
	}
}

func (l *logger) UnZap() *zap.Logger {
	return l.Logger
}

func (l *logger) Zap(z *zap.Logger) {
	l.Logger = z
}
