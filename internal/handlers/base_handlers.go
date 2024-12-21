package handlers

import (
	"io"
	"net/http"
	"shorter/internal/logger"
	"shorter/internal/storage"
	"shorter/pkg"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	Config  *Config
	Storage *storage.ShortURLStorage
	Log     logger.MyLogger
}

func NewHandler(config *Config, s *storage.ShortURLStorage, log logger.MyLogger) *Handler {
	return &Handler{
		Config:  config,
		Storage: s,
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
	id, err := h.Storage.Save(url)
	if err != nil {
		h.Log.Error("Ошибка сохранения URL: ", zap.Error(err))
	}
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
