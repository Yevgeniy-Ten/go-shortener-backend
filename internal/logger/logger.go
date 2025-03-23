// Description: Logger package for logging logic with context, using zap logger
package logger

import (
	"context"

	"go.uber.org/zap"
)

type contextKey struct{}

var zapFieldsKey = contextKey{}

// ZapFields is a slice of zap fields for context
type ZapFields []zap.Field

// ZapLogger base struct for logging
type ZapLogger struct {
	Log *zap.Logger
}

// Append appends fields to the zap fields
func (z *ZapFields) Append(fields ...zap.Field) {
	*z = append(*z, fields...)
}

// WithContextFields adds fields to the context for log with history
func (z *ZapLogger) WithContextFields(ctx context.Context, fields ...zap.Field) context.Context {
	return context.WithValue(ctx, zapFieldsKey, fields)
}
func (z *ZapLogger) withCtxFields(ctx context.Context, fields ...zap.Field) []zap.Field {
	ctxFields, ok := ctx.Value(zapFieldsKey).(ZapFields)
	if ok {
		ctxFields.Append(fields...)
	} else {
		ctxFields = make(ZapFields, 0, len(fields))
		ctxFields.Append(fields...)
	}
	return ctxFields
}

// InfoCtx logs info message
func (z *ZapLogger) InfoCtx(ctx context.Context, msg string, fields ...zap.Field) {
	z.Log.Info(msg, z.withCtxFields(ctx, fields...)...)
}

// ErrorCtx logs error message
func (z *ZapLogger) ErrorCtx(ctx context.Context, msg string, fields ...zap.Field) {
	z.Log.Error(msg, z.withCtxFields(ctx, fields...)...)
}

// InitLogger initializes logger
func InitLogger() (*ZapLogger, error) {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	myLogger, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	//nolint:errcheck // ignore error because it's not important
	defer myLogger.Sync()
	return &ZapLogger{
		myLogger,
	}, nil
}
