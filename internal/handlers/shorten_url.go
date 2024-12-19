package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"shorter/internal/storage"
	"shorter/pkg"
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
	//if err := c.BindJSON(&req); err != nil {
	//	c.String(http.StatusBadRequest, "Ошибка чтения тела запроса.")
	//	return
	//}

	if !pkg.ValidateURL(data.URL) {
		c.String(http.StatusBadRequest, "Некорректный URL.")
		return
	}
	id := storage.GlobalURLStorage.Save(data.URL)

	var responseData = ShortenResponse{
		Result: "http://localhost:8080/" + id,
	}
	var buf bytes.Buffer

	err = json.NewDecoder(&buf).Decode(&responseData)
	if err != nil {
		c.String(http.StatusInternalServerError, "Ошибка сериализации данных")
		return
	}
	c.JSON(http.StatusCreated, responseData)
}
