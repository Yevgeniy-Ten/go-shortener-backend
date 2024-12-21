package handlers

import (
	"io"
	"net/http"
	"shorter/internal/gzipper"
	"shorter/internal/logger"
	"shorter/internal/storage"
	"shorter/pkg"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Config *Config
}

func (h *Handler) PostHandler(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "Ошибка чтения тела запроса.")
		return
	}
	url := string(body)
	if !pkg.ValidateURL(url) {
		c.String(http.StatusBadRequest, "Некорректный URL.")
		return
	}
	id := storage.GlobalURLStorage.Save(url)
	respText := h.Config.ServerAddr + "/" + id
	c.String(http.StatusCreated, respText)
}
func (h *Handler) GetHandler(c *gin.Context) {
	id := c.Param("id")
	url := storage.GlobalURLStorage.GetURL(id)
	if url == "" {
		c.String(http.StatusBadRequest, "Not found")
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *Handler) CreateRouter() *gin.Engine {
	r := gin.Default()
	r.Use(gzipper.RequestResponseGzipMiddleware())
	r.Use(logger.RequestResponseInfoMiddleware())
	r.POST("/", h.PostHandler)
	r.POST("/api/shorten", h.ShortenURLHandler)
	r.GET("/:id", h.GetHandler)
	return r
}
