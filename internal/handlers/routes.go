package handlers

import (
	"context"
	"shorter/internal/gzipper"
	"shorter/internal/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handler) CreateRouter(
	middlewares ...gin.HandlerFunc,
) *gin.Engine {
	r := gin.Default()
	r.Use(middlewares...)
	r.POST("/", h.PostHandler)
	r.POST("/api/shorten", h.ShortenURLHandler)
	r.GET("/:id", h.GetHandler)
	return r
}

func (h *Handler) GetRoutes() *gin.Engine {
	ctx := context.Background()
	ctx = h.Log.WithContextFields(ctx,
		zap.String("Middleware", "RequestResponseInfoMiddleware"),
	)
	r := h.CreateRouter(gzipper.RequestResponseGzipMiddleware(),
		logger.RequestResponseInfoMiddleware(ctx, h.Log))
	return r
}
