package handlers

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"shorter/internal/app/lib"
	"shorter/internal/app/storage"
	"strings"
)

func PostHandler(c *gin.Context) {
	//contentType := req.Header.Get("Content-Type")
	contentType := c.GetHeader("Content-Type")
	if !strings.Contains(contentType, "text/plain") {
		c.String(http.StatusBadRequest, "Некорректный Content-Type.")
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "Ошибка чтения тела запроса.")
		return
	}
	url := string(body)
	if !lib.ValidateURL(url) {
		c.String(http.StatusBadRequest, "Некорректный URL.")
		return
	}
	id := storage.GlobalURLStorage.Save(url)
	respText := "http://localhost:8080/" + id
	c.String(http.StatusCreated, respText)
}
func GetHandler(c *gin.Context) {
	id := c.Param("id")
	url := storage.GlobalURLStorage.GetURL(id)
	if url == "" {
		c.String(http.StatusBadRequest, "Not found")
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, url)
}
