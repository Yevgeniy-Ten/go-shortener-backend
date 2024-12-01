package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"shorter/config"
	"shorter/internal/handlers"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {

	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}
	h := &handlers.Handler{
		Config: cfg.Config,
	}
	r := h.CreateRouter()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	fmt.Println("Starting server at", cfg.Address)
	return r.Run(cfg.Address)
}
