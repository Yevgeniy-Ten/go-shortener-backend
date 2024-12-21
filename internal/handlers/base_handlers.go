package handlers

import (
	"io"
	"net/http"
	"shorter/internal/logger"
	"shorter/internal/storage"
	"shorter/pkg"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Config  *Config
	Storage *storage.ShortURLStorage
	Log     logger.MyLogger
}

func NewHandler(config *Config, storage *storage.ShortURLStorage, log logger.MyLogger) *Handler {
	return &Handler{
		Config:  config,
		Storage: storage,
		Log:     log,
	}
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
	id := h.Storage.Save(url)
	respText := h.Config.ServerAddr + "/" + id
	c.String(http.StatusCreated, respText)
}

func (h *Handler) GetHandler(c *gin.Context) {
	id := c.Param("id")
	url := h.Storage.GetURL(id)
	if url == "" {
		c.String(http.StatusBadRequest, "Not found")
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, url)
}
