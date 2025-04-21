package logger

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RequestResponseInfoMiddleware logs request and response info
func RequestResponseInfoMiddleware(
	ctx context.Context,
	logger *ZapLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		uri := c.Request.RequestURI
		method := c.Request.Method
		c.Next()
		latency := time.Since(t)
		statusCodeToSent := c.Writer.Status()
		bodySizeToSent := c.Writer.Size()
		logger.InfoCtx(ctx, "RequestResponseInfo", zap.String("uri", uri), zap.String("method", method),
			zap.Int("status", statusCodeToSent), zap.Int("size", bodySizeToSent), zap.Duration("latency", latency))
	}
}
