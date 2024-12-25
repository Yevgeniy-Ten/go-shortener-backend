package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"shorter/internal/domain"
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
		c.String(http.StatusBadRequest, "Некорректный URL.")
		return
	}
	id, err := h.Storage.Save(data.URL)
	if err != nil {
		h.Log.ErrorCtx(context.TODO(), "Error when save ", zap.Error(err))
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
