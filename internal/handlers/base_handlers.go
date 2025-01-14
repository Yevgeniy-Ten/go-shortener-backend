package handlers

import (
	"errors"
	"io"
	"net/http"
	"shorter/internal/storage/database/urls"
	"shorter/pkg"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handler) PostHandler(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "Read error")
		return
	}
	url := string(body)
	if !pkg.ValidateURL(url) {
		c.String(http.StatusBadRequest, "Некорректный ShortURL.")
		return
	}
	id, err := h.Storage.Save(url)
	if err != nil {
		var duplicateError *urls.DuplicateError
		if errors.As(err, &duplicateError) {
			c.String(http.StatusConflict, h.Config.ServerAddr+"/"+duplicateError.ShortURL)
			return
		}
		h.Log.Log.Error("Error when save", zap.Error(err))
		c.String(http.StatusInternalServerError, "Error")
		return
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
