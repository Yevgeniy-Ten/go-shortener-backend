package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shorter/internal/gzipper"
	"shorter/internal/logger"
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
	r := h.CreateRouter(gzipper.RequestResponseGzipMiddleware(),
		logger.RequestResponseInfoMiddleware(h.Log))
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	return r
}
