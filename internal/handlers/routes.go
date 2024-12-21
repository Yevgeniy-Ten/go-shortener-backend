package handlers

import (
	"shorter/internal/gzipper"
	"shorter/internal/logger"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateRouter() *gin.Engine {
	r := gin.Default()
	r.Use(gzipper.RequestResponseGzipMiddleware(h.Log))
	r.Use(logger.RequestResponseInfoMiddleware(h.Log))
	r.POST("/", h.PostHandler)
	r.POST("/api/shorten", h.ShortenURLHandler)
	r.GET("/:id", h.GetHandler)
	return r
}
