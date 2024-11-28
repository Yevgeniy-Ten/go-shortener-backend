package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"shorter/internal/handlers"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {

	cfg := NewConfig()
	h := &handlers.Handler{
		Config: &cfg.Config,
	}
	r := h.CreateRouter()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	return r.Run(":" + *cfg.port)
}
