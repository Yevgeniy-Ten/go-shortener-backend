package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"shorter/pkg"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	Result string `json:"result"`
}

func (h *Handler) ShortenURLHandler(c *gin.Context) {
	var data ShortenRequest
	body, err := c.GetRawData()

	if err != nil {
		c.String(http.StatusBadRequest, "Ошибка чтения тела запроса.")
		return
	}

	if err := json.Unmarshal(body, &data); err != nil {
		c.String(http.StatusBadRequest, "Ошибка чтения тела запроса.")
		return
	}

	if !pkg.ValidateURL(data.URL) {
		c.String(http.StatusBadRequest, "Некорректный URL.")
		return
	}
	id, err := h.Storage.Save(data.URL)
	if err != nil {
		h.Log.ErrorCtx(context.TODO(), "Ошибка сохранения URL: ", zap.Error(err))
	}
	var responseData = ShortenResponse{
		Result: "http://localhost:8080/" + id,
	}
	var buf bytes.Buffer

	err = json.NewEncoder(&buf).Encode(&responseData)
	if err != nil {
		c.String(http.StatusInternalServerError, "Ошибка сериализации данных")
		return
	}
	c.JSON(http.StatusCreated, responseData)
}
