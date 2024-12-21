package main

import (
	"log"
	"net/http"
	"shorter/config"
	"shorter/internal/gzipper"
	"shorter/internal/handlers"
	"shorter/internal/logger"
	"shorter/internal/storage"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
	myLogger, err := logger.InitLogger()
	if err != nil {
		return err
	}
	fileStorage, err := storage.NewShortURLStorage(cfg.FilePath)
	if err != nil {
		return err
	}
	defer fileStorage.Close()
	h := handlers.NewHandler(cfg.Config, fileStorage, myLogger)
	r := h.CreateRouter(
		gzipper.RequestResponseGzipMiddleware(),
		logger.RequestResponseInfoMiddleware(h.Log),
	)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	myLogger.Info("Server started", zap.String("address", cfg.Address))
	return r.Run(cfg.Address)
}
