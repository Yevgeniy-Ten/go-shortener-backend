package logger

import (
	"context"

	"go.uber.org/zap"
)

type contextKey struct{}

// Создайте уникальный экземпляр ключа
var zapFieldsKey = contextKey{}

type ZapFields []zap.Field
type ZapLogger struct {
	Log *zap.Logger
}

func (z *ZapFields) Append(fields ...zap.Field) {
	*z = append(*z, fields...)
}

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

func (z *ZapLogger) InfoCtx(ctx context.Context, msg string, fields ...zap.Field) {
	z.Log.Info(msg, z.withCtxFields(ctx, fields...)...)
}
func (z *ZapLogger) ErrorCtx(ctx context.Context, msg string, fields ...zap.Field) {
	z.Log.Error(msg, z.withCtxFields(ctx, fields...)...)
}
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
