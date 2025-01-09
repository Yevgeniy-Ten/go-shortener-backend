package handlers

import (
	"context"
	"io"
	"net/http"
	"shorter/internal/logger"
	"shorter/pkg"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type storage interface {
	Save(value string) (string, error)
	GetURL(shortURL string) string
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

func (h *Handler) PostHandler(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "Read error")
		return
	}
	url := string(body)
	if !pkg.ValidateURL(url) {
		c.String(http.StatusBadRequest, "Некорректный URL.")
		return
	}
	id, err := h.Storage.Save(url)
	if err != nil {
		h.Log.ErrorCtx(context.TODO(), "Error when save ", zap.Error(err))
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
