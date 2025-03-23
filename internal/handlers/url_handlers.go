package handlers

import (
	"errors"
	"io"
	"net/http"
	"shorter/internal/cookies"
	"shorter/internal/urlstorage/database/urls"
	"shorter/pkg"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handler) PostHandler(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Data(http.StatusBadRequest, "text/plain", []byte("Read error"))
		return
	}
	url := string(body)
	if !pkg.ValidateURL(url) {
		c.Data(http.StatusBadRequest, "text/plain", []byte("Incorrect ShortURL"))
		return
	}
	var (
		urlID  string
		userID int
	)
	if h.Storage.User != nil {
		if userID, err = cookies.GetUserFromCookie(c); err != nil {
			c.Data(http.StatusUnauthorized, "text/plain", []byte("Unauthorized"))
			return
		}
	}
	urlID, err = h.Storage.URLS.Save(url, userID)
	baseURL := h.Config.ServerAddr + "/"
	if err != nil {
		var duplicateError *urls.DuplicateError
		if errors.As(err, &duplicateError) {
			c.String(http.StatusConflict, baseURL+duplicateError.ShortURL)
			return
		}
		h.l.Log.Error("Error when save", zap.Error(err))
		c.String(http.StatusInternalServerError, "Error")
		return
	}
	var sb strings.Builder
	sb.WriteString(baseURL)
	sb.WriteString(urlID)
	c.Data(http.StatusCreated, "text/plain", []byte(sb.String()))
}

func (h *Handler) GetHandler(c *gin.Context) {
	id := c.Param("id")
	url, err := h.Storage.URLS.GetURL(id)
	if err != nil {
		var urlIsDeletedError *urls.URLIsDeletedError
		if errors.As(err, &urlIsDeletedError) {
			c.Data(http.StatusGone, "text/plain", []byte("URL is deleted"))
			return
		}
		h.l.Log.Error("Error when get", zap.Error(err))
		c.Data(http.StatusNotFound, "text/plain", []byte("Not found"))
		return
	}
	if url == "" {
		c.Status(http.StatusBadRequest)
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, url)
}
func (h *Handler) GetUserUrls(c *gin.Context) {
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
