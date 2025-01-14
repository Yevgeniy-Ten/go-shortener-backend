package handlers

import (
	"context"
	"net/http"
	"shorter/internal/domain"
	"shorter/internal/gzipper"
	"shorter/internal/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type storage interface {
	Save(value string) (string, error)
	GetURL(shortURL string) string
	SaveBatch(urls []domain.URLS) error
	Ping() error
}

type Handler struct {
	Config  *Config
	Storage storage
	Log     *logger.ZapLogger
}

func NewHandler(config *Config, s storage, log *logger.ZapLogger) *Handler {
	return &Handler{
		Config:  config,
		Storage: s,
		Log:     log,
	}
}

func InitHandlers(config *Config, s storage, log *logger.ZapLogger) *gin.Engine {
	ctx := context.Background()
	h := NewHandler(config, s, log)
	ctx = h.Log.WithContextFields(ctx,
		zap.String("Middleware", "RequestResponseInfoMiddleware"),
	)
	r := h.CreateRouter(gzipper.RequestResponseGzipMiddleware(),
		logger.RequestResponseInfoMiddleware(ctx, h.Log))
	r.GET("/ping", func(c *gin.Context) {
		if err := h.Storage.Ping(); err != nil {
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
