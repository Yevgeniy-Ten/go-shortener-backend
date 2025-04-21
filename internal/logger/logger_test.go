package logger

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

// Test WithContextFields
func TestWithContextFields(t *testing.T) {
	a := assert.New(t)
	logger, _ := InitLogger()
	ctx := context.Background()

	field1 := zap.String("key1", "value1")
	field2 := zap.Int("key2", 42)

	ctx = logger.WithContextFields(ctx, field1)
	ctx = logger.WithContextFields(ctx, field2)

	fields, ok := ctx.Value(zapFieldsKey).([]zap.Field)
	a.True(ok, "Expected zap fields in context")
	a.Len(fields, 1, "Expected only the last added field to be stored")
	a.Equal(field2, fields[0], "Expected the last added field to be present")
}

// Test InfoCtx and ErrorCtx logging
func TestLoggingWithContext(t *testing.T) {
	a := assert.New(t)
	core := zaptest.NewLogger(t)
	logger := &ZapLogger{Log: core}
	ctx := context.Background()

	field := zap.String("testKey", "testValue")
	ctx = logger.WithContextFields(ctx, field)

	// Здесь можно дополнительно замокать логгер, если нужно проверять записи
	a.NotNil(logger, "Logger should not be nil")
	logger.InfoCtx(ctx, "test message")
	logger.ErrorCtx(ctx, "error message")
}

// Test InitLogger
func TestInitLogger(t *testing.T) {
	a := assert.New(t)
	logger, err := InitLogger()

	a.NoError(err, "Logger should initialize without errors")
	a.NotNil(logger.Log, "Logger instance should not be nil")
}
