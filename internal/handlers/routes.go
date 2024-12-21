package handlers

import (
	"github.com/gin-gonic/gin"
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
