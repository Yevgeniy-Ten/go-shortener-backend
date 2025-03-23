package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"shorter/internal/cookies"
	"shorter/internal/domain"
	"shorter/internal/urlstorage/database/urls"
	"shorter/pkg"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ShortenResponse struct {
	Result string `json:"result"`
}

func (h *Handler) ShortenURLHandler(c *gin.Context) {
	var data domain.ShortenRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&data); err != nil {
		c.String(http.StatusBadRequest, "Read error")
		return
	}

	if !pkg.ValidateURL(data.URL) {
		c.String(http.StatusBadRequest, "Некорректный ShortURL.")
		return
	}
	var (
		urlID  string
		userID int
		err    error
	)
	if h.Storage.User != nil {
		if userID, err = cookies.GetUserFromCookie(c); err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			return
		}
	}
	urlID, err = h.Storage.URLS.Save(data.URL, userID)
	if err != nil {
		var duplicateError *urls.DuplicateError
		if errors.As(err, &duplicateError) {
			c.JSON(http.StatusConflict, ShortenResponse{
				Result: h.Config.ServerAddr + "/" + duplicateError.ShortURL,
			})
			return
		}
		c.String(http.StatusInternalServerError, "Error")
		return
	}
	var sb strings.Builder
	sb.WriteString(h.Config.ServerAddr)
	sb.WriteString("/")
	sb.WriteString(urlID)
	responseData := ShortenResponse{Result: sb.String()}
	c.JSON(http.StatusCreated, responseData)
}

func (h *Handler) ShortenURLSHandler(c *gin.Context) {
	var data []domain.URLS
	if err := json.NewDecoder(c.Request.Body).Decode(&data); err != nil {
		c.String(http.StatusBadRequest, "Read error")
		return
	}
	var userID int
	var err error
	if h.Storage.User != nil {
		if userID, err = cookies.GetUserFromCookie(c); err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			return
		}
	}
	err = h.Storage.URLS.SaveBatch(data, userID)
	if err != nil {
		h.l.Log.Error("Error when save ", zap.Error(err))
		c.String(http.StatusInternalServerError, "Error")
	}
	responseData := make([]domain.ShortenerBatchResponse, 0, len(data))

	for _, url := range data {
		var sb strings.Builder
		sb.WriteString(h.Config.ServerAddr)
		sb.WriteString("/")
		sb.WriteString(url.CorrelationID)
		responseData = append(responseData, domain.ShortenerBatchResponse{
			CorrelationID: url.CorrelationID,
			ShortURL:      sb.String(),
		})
	}
	c.JSON(http.StatusCreated, responseData)
}
func (h *Handler) DeleteMyUrls(c *gin.Context) {
	var userID, err = cookies.GetUserFromCookie(c)
	if err != nil {
		c.String(http.StatusUnauthorized, "Unauthorized")
		return
	}
	var correlationIDS []string
	if err := c.ShouldBindJSON(&correlationIDS); err != nil {
		c.String(http.StatusBadRequest, "Bad request")
		return
	}
	go func() {
		err := h.Storage.URLS.DeleteURLs(correlationIDS, userID)
		if err != nil {
			h.l.Log.Error("Error when delete urls", zap.Error(err))
		}
	}()
	c.Status(http.StatusAccepted)
}
