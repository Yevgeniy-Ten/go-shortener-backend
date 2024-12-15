package logger

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

func RequestResponseInfoMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		uri := c.Request.RequestURI
		method := c.Request.Method
		c.Next()
		latency := time.Since(t)
		statusCodeToSent := c.Writer.Status()
		bodySizeToSent := c.Writer.Size()
		Log.Info("Request", zap.String("uri", uri), zap.String("method", method), zap.Int("status", statusCodeToSent), zap.Int("size", bodySizeToSent), zap.Duration("latency", latency))
	}
}
