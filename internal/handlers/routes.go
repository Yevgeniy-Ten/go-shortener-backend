package handlers

import (
	"github.com/gin-gonic/gin"
	"shorter/internal/gzipper"
	"shorter/internal/logger"
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
