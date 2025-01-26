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
	l       *logger.ZapLogger
}

func NewHandler(config *Config, s domain.Storage, log *logger.ZapLogger) *Handler {
	return &Handler{
		Config:  config,
		Storage: s,
		l:       log,
	}
}

func InitHandlers(config *Config, s domain.Storage, log *logger.ZapLogger, withDatabase bool) *gin.Engine {
	ctx := context.Background()
	h := NewHandler(config, s, log)
	ctx = h.l.WithContextFields(ctx,
		zap.String("Middleware", "RequestResponseInfoMiddleware"),
	)

	middlewares := []gin.HandlerFunc{
		gzipper.RequestResponseGzipMiddleware(),
		logger.RequestResponseInfoMiddleware(ctx, h.l),
	}
	r := h.CreateRouter(withDatabase, middlewares...)
	r.GET("/ping", func(c *gin.Context) {
		if !withDatabase {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Status(http.StatusOK)
	})
	return r
}

func (h *Handler) CreateRouter(withDatabase bool,
	middlewares ...gin.HandlerFunc,
) *gin.Engine {
	r := gin.Default()
	r.Use(middlewares...)
	r.POST("/", cookies.CreateUserMiddleware(withDatabase, h.l, h.Storage.User), h.PostHandler)
	r.POST("/api/shorten", cookies.CreateUserMiddleware(withDatabase, h.l, h.Storage.User), h.ShortenURLHandler)
	r.POST("/api/shorten/batch", cookies.CreateUserMiddleware(withDatabase, h.l, h.Storage.User), h.ShortenURLSHandler)
	r.GET("/:id", h.GetHandler)
	r.GET("/api/user/urls", h.GetAllMyUrls)
	r.DELETE("/api/user/urls", h.DeleteMyUrls)
	return r
}
