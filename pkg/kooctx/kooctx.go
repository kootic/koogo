package kooctx

import (
	"context"

	"go.uber.org/zap"
)

type contextKey string

const (
	ContextKeyLogger contextKey = "logger"
)

func getValueFromContext[T any](ctx context.Context, key contextKey) (T, bool) {
	value, ok := ctx.Value(key).(T)
	return value, ok
}

func SetContextLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, ContextKeyLogger, logger)
}

func GetContextLogger(ctx context.Context) *zap.Logger {
	logger, ok := getValueFromContext[*zap.Logger](ctx, ContextKeyLogger)
	if !ok {
		return zap.NewNop()
	}

	return logger
}

func WithLoggerFields(ctx context.Context, fields ...zap.Field) (context.Context, *zap.Logger) {
	logger := GetContextLogger(ctx)
	newLogger := logger.With(fields...)

	return SetContextLogger(ctx, newLogger), newLogger
}
