package logger

import (
	"context"

	"go.uber.org/zap"
)

const zapFieldsKey = "zapFields"

type ZapFields []zap.Field
type ZapLogger struct {
	log *zap.Logger
}

func (z *ZapFields) Append(fields ...zap.Field) {
	*z = append(*z, fields...)
}

func (z *ZapLogger) WithContextFields(ctx context.Context, fields ...zap.Field) context.Context {
	//nolint:staticcheck // ignore SA1019 because we guarantee that zapFieldsKey is unique
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
	z.log.Info(msg, z.withCtxFields(ctx, fields...)...)
}
func (z *ZapLogger) ErrorCtx(ctx context.Context, msg string, fields ...zap.Field) {
	z.log.Error(msg, z.withCtxFields(ctx, fields...)...)
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
