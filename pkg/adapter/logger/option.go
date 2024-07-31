package logger

import "context"

type Option func(*logger)

type ContextEncoder func(context.Context, Logger) Logger

func WithContextEncoder(contextEncoder ContextEncoder) Option {
	return func(l *logger) {
		l.contextEncoder = contextEncoder
	}
}
