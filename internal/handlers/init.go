package handlers

import (
	"context"
	"net/http"
	"shorter/internal/cookies"
	"shorter/internal/domain"
	"shorter/internal/gzipper"
	"shorter/internal/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	Config  *Config
	Storage domain.Storage
	Log     *logger.ZapLogger
}

func NewHandler(config *Config, s domain.Storage, log *logger.ZapLogger) *Handler {
	return &Handler{
		Config:  config,
		Storage: s,
		Log:     log,
	}
}

func InitHandlers(config *Config, s domain.Storage, log *logger.ZapLogger, withDatabase bool) *gin.Engine {
	ctx := context.Background()
	h := NewHandler(config, s, log)
	ctx = h.Log.WithContextFields(ctx,
		zap.String("Middleware", "RequestResponseInfoMiddleware"),
	)

	middlewares := []gin.HandlerFunc{
		gzipper.RequestResponseGzipMiddleware(),
		logger.RequestResponseInfoMiddleware(ctx, h.Log),
	}
	if withDatabase {
		middlewares = append(middlewares, cookies.CreateUserMiddleware(h.Log, h.Storage.User))
	}
	r := h.CreateRouter(middlewares...)
	r.GET("/ping", func(c *gin.Context) {
		if !withDatabase {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Status(http.StatusOK)
	})
	return r
}

func (h *Handler) CreateRouter(
	middlewares ...gin.HandlerFunc,
) *gin.Engine {
	r := gin.Default()
	r.Use(middlewares...)
	r.POST("/", h.PostHandler)
	r.POST("/api/shorten", h.ShortenURLHandler)
	r.POST("/api/shorten/batch", h.ShortenURLSHandler)
	r.GET("/:id", h.GetHandler)

	return r
}
