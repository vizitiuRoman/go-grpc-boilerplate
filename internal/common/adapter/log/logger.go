package log

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

	Named(string) Logger
	WithComponent(context.Context, string) Logger
	WithMethod(context.Context, string) Logger

	WithCtx(ctx context.Context, fields ...zap.Field) Logger
	With(fields ...zap.Field) Logger
	Zap() *zap.Logger
}

type logger struct {
	*zap.Config
	*zap.Logger
}

func NewLogger(cfg *zap.Config, fields ...zap.Field) (Logger, error) {
	baseLogger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	baseLogger = baseLogger.With(fields...)

	defer func() { _ = baseLogger.Sync() }()

	return &logger{
		Config: cfg,
		Logger: baseLogger,
	}, nil
}

func defaultConsoleConfig(lvl string) *zap.Config {
	zap.NewProductionConfig()
	return &zap.Config{
		Level:            unmarshalLevel(lvl),
		DisableCaller:    false,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
		Encoding:         "console",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "msg",
			LevelKey:       "level",
			TimeKey:        "ts",
			NameKey:        "logger",
			CallerKey:      "file",
			StacktraceKey:  "stacktrace",
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,
		},
	}
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
			TimeKey:        "ts",
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

func MustDefaultJSONLogger(lvl string) Logger {
	l, err := NewLogger(defaultJSONConfig(lvl))
	if err != nil {
		panic(err)
	}

	return l
}

func MustDefaultConsoleLogger(lvl string) Logger {
	l, err := NewLogger(defaultConsoleConfig(lvl))
	if err != nil {
		panic(err)
	}

	return l
}

func unmarshalLevel(lvl string) zap.AtomicLevel {
	level := zap.NewAtomicLevel()
	err := level.UnmarshalText([]byte(lvl))
	if err != nil || lvl == "" {
		level.SetLevel(zap.DebugLevel)
	}

	return level
}

func (l *logger) Named(name string) Logger {
	return &logger{
		Config: l.Config,
		Logger: l.Logger.Named(name),
	}
}

func (l *logger) WithComponent(ctx context.Context, componentName string) Logger {
	const key = "go.component"

	return l.WithCtx(ctx, zap.String(key, componentName))
}

func (l *logger) WithMethod(ctx context.Context, methodName string) Logger {
	const key = "go.method"

	return l.WithCtx(ctx, zap.String(key, methodName))
}

func (l *logger) With(fields ...zap.Field) Logger {
	return &logger{
		Config: l.Config,
		Logger: l.Logger.With(fields...),
	}
}

func (l *logger) WithCtx(ctx context.Context, fields ...zap.Field) Logger {
	// Extract values info from context: span, metadata, session information, etc.
	// Add required fields to logger

	res := &logger{
		Config: l.Config,
		Logger: l.Logger.With(fields...),
	}

	return res
}

func (l *logger) Zap() *zap.Logger {
	return l.Logger
}
