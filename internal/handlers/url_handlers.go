package handlers

import (
	"errors"
	"io"
	"net/http"
	"shorter/internal/cookies"
	"shorter/internal/urlstorage/database/urls"
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
	var urlID string
	var userID int
	if h.Storage.User != nil {
		if userID, err = cookies.GetUserFromCookie(c); err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			return
		}
	}
	urlID, err = h.Storage.URLS.Save(url, userID)
	if err != nil {
		var duplicateError *urls.DuplicateError
		if errors.As(err, &duplicateError) {
			c.String(http.StatusConflict, h.Config.ServerAddr+"/"+duplicateError.ShortURL)
			return
		}
		h.l.Log.Error("Error when save", zap.Error(err))
		c.String(http.StatusInternalServerError, "Error")
		return
	}
	respText := h.Config.ServerAddr + "/" + urlID
	c.String(http.StatusCreated, respText)
}

func (h *Handler) GetHandler(c *gin.Context) {
	id := c.Param("id")
	url, err := h.Storage.URLS.GetURL(id)
	if err != nil {
		var urlIsDeletedError *urls.URLIsDeletedError
		if errors.As(err, &urlIsDeletedError) {
			c.String(http.StatusGone, "URL is deleted")
			return
		}
		h.l.Log.Error("Error when get", zap.Error(err))
		c.String(http.StatusNotFound, "Not found")
		return
	}
	if url == "" {
		c.String(http.StatusBadRequest, "Not found")
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, url)
}
func (h *Handler) GetAllMyUrls(c *gin.Context) {
	var userID int
	var err error
	if userID, err = cookies.GetUserFromCookie(c); err != nil {
		c.String(http.StatusUnauthorized, "Unauthorized")
		return
	}
	userUrls, err := h.Storage.URLS.GetUserURLs(userID, h.Config.ServerAddr+"/")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error")
		return
	}
	if len(userUrls) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, userUrls)
}
