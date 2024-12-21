package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"shorter/pkg"

	"github.com/gin-gonic/gin"
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

	// #second variant
	//  if err := c.BindJSON(&req); err != nil {
	//	c.String(http.StatusBadRequest, "Ошибка чтения тела запроса.")
	//	return
	// }

	if !pkg.ValidateURL(data.URL) {
		c.String(http.StatusBadRequest, "Некорректный URL.")
		return
	}
	id := h.Storage.Save(data.URL)
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
