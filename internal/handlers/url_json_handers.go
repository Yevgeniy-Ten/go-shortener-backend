package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"shorter/internal/cookies"
	"shorter/internal/domain"
	"shorter/internal/urlstorage/database/urls"
	"shorter/pkg"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ShortenResponse struct {
	Result string `json:"result"`
}

func (h *Handler) ShortenURLHandler(c *gin.Context) {
	var data domain.ShortenRequest
	body, err := c.GetRawData()

	if err != nil {
		c.String(http.StatusBadRequest, "Read error")
		return
	}
	if err := json.Unmarshal(body, &data); err != nil {
		c.String(http.StatusBadRequest, "Read error")
		return
	}

	if !pkg.ValidateURL(data.URL) {
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
	var responseData = ShortenResponse{
		Result: h.Config.ServerAddr + "/" + urlID,
	}
	var buf bytes.Buffer

	err = json.NewEncoder(&buf).Encode(&responseData)
	if err != nil {
		c.String(http.StatusInternalServerError, "Ошибка сериализации данных")
		return
	}
	c.JSON(http.StatusCreated, responseData)
}

func (h *Handler) ShortenURLSHandler(c *gin.Context) {
	var data []domain.URLS
	body, err := c.GetRawData()

	if err != nil {
		c.String(http.StatusBadRequest, "Read error")
		return
	}

	if err := json.Unmarshal(body, &data); err != nil {
		c.String(http.StatusBadRequest, "Read error")
		return
	}
	var userID int
	if h.Storage.User != nil {
		if userID, err = cookies.GetUserFromCookie(c); err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			return
		}
	}
	err = h.Storage.URLS.SaveBatch(data, userID)
	if err != nil {
		h.Log.Log.Error("Error when save ", zap.Error(err))
		c.String(http.StatusInternalServerError, "Error")
	}
	var responseData []domain.ShortenerBatchResponse
	for _, url := range data {
		responseData = append(responseData, domain.ShortenerBatchResponse{
			CorrelationID: url.CorrelationID,
			ShortURL:      h.Config.ServerAddr + "/" + url.CorrelationID,
		})
	}
	c.JSON(http.StatusCreated, responseData)
}
